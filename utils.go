package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func makeWazuhRequest(client *http.Client, req *http.Request, headers map[string]string, timeout int) (map[string]interface{}, error) {
	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, http.ErrHandlerTimeout) {
			return nil, fmt.Errorf("connection to manager timed out after %d seconds", timeout)
		}
		return nil, fmt.Errorf("error connecting to manager: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed: %s", resp.Status)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response from manager: %s", resp.Status)
	}

	// Not using ioutil because it is deprecated
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
