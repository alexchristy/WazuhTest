package main

// main
func main() {

	args := parseArguments()

	// Initialize the WazuhServer object
	wazuhServer, err := NewWazuhServer(args.User, args.Password, args.Host, args.Timeout)
	if err != nil {
		PrintRed("Error initializing WazuhServer object: " + err.Error())
		return
	}

	wazuhServer.checkConnection(args.Verbosity)

	runTestGroup(*wazuhServer, "./tests", args.Threads, args.Verbosity, args.Timeout)
}
