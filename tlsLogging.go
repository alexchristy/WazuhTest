package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

func enableTLSKeyLogging(client *http.Client, logFilePath string) error {
	// Open the log file
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	// Check if the client already has an existing transport
	if client.Transport != nil {
		transport, ok := client.Transport.(*http.Transport)
		if !ok {
			return fmt.Errorf("passed in client does not have a valid  *http.Transport")
		}

		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}

		transport.TLSClientConfig.KeyLogWriter = file

	} else {
		// Create new transport
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				KeyLogWriter: file,
			},
		}
	}

	return nil
}
