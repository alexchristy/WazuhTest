package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
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
	token         string
	protocol      string
	port          int
	loginEndpoint string
	httpClient    *http.Client
}

func NewWazuhServer(ApiUser string, ApiPass string, Hostname string, Timeout int) (*WazuhServer, error) {
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
	ws.httpClient = &http.Client{
		Timeout: time.Duration(ws.Timeout) * time.Second,
		// Do not verify certs
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Attempt to authenticate to the manager
	err := ws.getToken()
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (ws *WazuhServer) getToken() error {
	loginUrl := fmt.Sprintf("%s://%s:%d/%s", ws.protocol, ws.Hostname, ws.port, ws.loginEndpoint)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", ws.ApiUser, ws.ApiPass)))
	loginHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Basic %s", basicAuth),
	}

	PrintWhite("Authenticating to manager: " + ws.Hostname)

	req, err := http.NewRequest("POST", loginUrl, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	result, err := makeWazuhRequest(ws.httpClient, req, loginHeaders, ws.Timeout)
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

	fmt.Println()

	return nil
}

func (ws *WazuhServer) checkConnection(verbosity int) error {
	PrintWhite("Verifying connection to manager...")

	if ws.token == "" {
		return fmt.Errorf("no token available. Please authenticate to the manager first")
	}

	checkUrl := fmt.Sprintf("%s://%s:%d/", ws.protocol, ws.Hostname, ws.port)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", ws.token),
	}

	req, err := http.NewRequest("GET", checkUrl, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	result, err := makeWazuhRequest(ws.httpClient, req, headers, ws.Timeout)
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

	fmt.Println()

	return nil
}
