package main

import "os"

func cliExit(numFailedTests int) {
	// Exit with an error if at least one
	// test failed.
	if numFailedTests > 0 {
		os.Exit(1)
	}

	os.Exit(0)
}
