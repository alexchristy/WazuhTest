package main

// main
func main() {

	args := parseArguments()

	// Initialize the WazuhServer object
	wazuhServer, err := NewWazuhServer(args.User, args.Password, args.Host, args.Timeout, args.TlsLogPath)
	if err != nil {
		PrintRed("Error initializing WazuhServer object: " + err.Error())
		return
	}

	wazuhServer.checkConnection(args.Verbosity)

	numTests, numFailedTests, numWarnTests, err := runTestGroup(wazuhServer, args.TestsDir, args.Threads, args.Verbosity)
	if err != nil {
		panic(err)
	}

	printSummary(numTests, numFailedTests, numWarnTests)

	if args.CliMode {
		cliExit(numFailedTests)
	}
}
