package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type WazuhServer struct {
	ApiUser   string
	ApiPass   string
	Hostname  string
	Connected bool
	Timeout   int

	// Internal variables
	token           string
	protocol        string
	port            int
	loginEndpoint   string
	logTestEndpoint string
	httpClient      *http.Client
	sessionToken    string
}

func NewWazuhServer(ApiUser string, ApiPass string, Hostname string, Timeout int, tlsKeyLogPath string) (*WazuhServer, error) {
	ws := new(WazuhServer)

	// Validate the input
	if ApiUser == "" {
		panic("Api user cannot be empty")
	}

	if ApiPass == "" {
		panic("Api password cannot be empty")
	}

	if Hostname == "" {
		panic("Wazuh manager server Hostname cannot be empty")
	}

	if Timeout < 0 {
		panic("Timeout cannot be less than 0")
	}

	ws.ApiUser = ApiUser
	ws.ApiPass = ApiPass
	ws.Hostname = Hostname
	ws.Timeout = Timeout
	ws.protocol = "https"
	ws.port = 55000
	ws.loginEndpoint = "security/user/authenticate"
	ws.logTestEndpoint = "logtest"
	ws.httpClient = &http.Client{
		Timeout: time.Duration(ws.Timeout) * time.Second,
		// Do not verify certs
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	if len(tlsKeyLogPath) > 0 {
		enableTLSKeyLogging(ws.httpClient, tlsKeyLogPath)
	}

	// Attempt to authenticate to the manager
	err := ws.requestAuthToken()
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (ws *WazuhServer) requestAuthToken() error {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", ws.ApiUser, ws.ApiPass)))
	loginHeaders := map[string]interface{}{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Basic %s", basicAuth),
	}

	PrintWhite("Authenticating to manager: " + ws.Hostname)

	req, err := http.NewRequest("POST", ws.getLoginUrl(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	result, err := ws.sendRequest(req, loginHeaders)
	if err != nil {
		return err
	}

	// Response format:
	// {
	//   "data": {
	//     "token": "eyJhb..."
	//   }
	//   "error": 0
	// }
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected response format: no data field")
	}

	token, ok := data["token"].(string)
	if !ok {
		return fmt.Errorf("unexpected response format: no token field")
	}

	// Check if there error field is populated
	if result["error"] != float64(0) {
		return fmt.Errorf("error authenticating, manager reported an error")
	}

	ws.token = token

	PrintGreen("Sucessfully authenticated to manager.")

	return nil
}

func (ws *WazuhServer) getLoginUrl() string {
	return fmt.Sprintf("%s://%s:%d/%s", ws.protocol, ws.Hostname, ws.port, ws.loginEndpoint)
}

func (ws *WazuhServer) getLogTestUrl() string {
	return fmt.Sprintf("%s://%s:%d/%s", ws.protocol, ws.Hostname, ws.port, ws.logTestEndpoint)
}

func (ws *WazuhServer) getBaseUrl() string {
	return fmt.Sprintf("%s://%s:%d/", ws.protocol, ws.Hostname, ws.port)
}

func (ws *WazuhServer) getToken() string {
	return ws.token
}

func (ws *WazuhServer) hasSession() bool {
	return len(ws.sessionToken) > 0
}

func (ws *WazuhServer) getSessionToken() string {
	return ws.sessionToken
}

func (ws *WazuhServer) setSessionToken(token string) {
	ws.sessionToken = token
}

func (ws *WazuhServer) sendRequest(req *http.Request, headers map[string]interface{}) (map[string]interface{}, error) {
	// Add headers
	for key, value := range headers {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}

	// Check if the request has data
	if req.Body != nil {
		// Check if the Content-Type header is set
		if req.Header.Get("Content-Type") == "" {
			return nil, fmt.Errorf("Content-Type header is required when data is included in the request")
		}
	}

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, http.ErrHandlerTimeout) {
			return nil, fmt.Errorf("connection to manager timed out after %d seconds", ws.Timeout)
		}
		return nil, fmt.Errorf("error connecting to manager: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed: %s", resp.Status)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response from manager: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing token response body: %s", err)
	}

	return result, nil
}

func (ws *WazuhServer) checkConnection(verbosity int) error {
	PrintWhite("Verifying connection to manager...")

	if ws.token == "" {
		return fmt.Errorf("no token available. Please authenticate to the manager first")
	}

	headers := map[string]interface{}{
		"Authorization": fmt.Sprintf("Bearer %s", ws.token),
	}

	req, err := http.NewRequest("GET", ws.getBaseUrl(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	result, err := ws.sendRequest(req, headers) // data is empty
	if err != nil {
		return err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected response format: no data field")
	}

	title, ok := data["title"].(string)
	if !ok {
		return fmt.Errorf("unexpected response format: no title field")
	}

	if title != "Wazuh API REST" {
		return fmt.Errorf("bad response from manager. Returned title: %s", title)
	}

	apiVersion, ok := data["api_version"].(string)
	if !ok {
		return fmt.Errorf("unexpected response format: no api_version field")
	}

	revision, ok := data["revision"].(float64)
	if !ok {
		return fmt.Errorf("unexpected response format: no revision field")
	}

	if verbosity > 1 {
		PrintWhite("Wazuh API version: " + apiVersion + " (revision: " + strconv.FormatFloat(revision, 'f', -1, 64) + ")")
	} else {
		if verbosity > 0 {
			PrintWhite("Wazuh API version: " + apiVersion)
		}
	}
	PrintGreen("Verified connection to manager.")

	fmt.Printf("\n\n")

	return nil
}
