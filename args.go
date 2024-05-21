package main

import (
	"flag"
	"fmt"
	"os"
)

type Arguments struct {
	Host      string
	TestsDir  string
	User      string
	Password  string
	Threads   int
	Timeout   int
	Verbosity int
}

func parseArguments() Arguments {
	var args Arguments

	flag.StringVar(&args.TestsDir, "d", "./tests", "The directory containing the test groups. Defaults to './tests'.")
	flag.StringVar(&args.User, "u", "wazuh", "The username for the Wazuh API. Defaults to 'wazuh'.")
	flag.StringVar(&args.Password, "p", "wazuh", "The password for the Wazuh API. Defaults to 'wazuh'.")
	flag.IntVar(&args.Threads, "t", 1, "The number of threads to use for running tests. Defaults to 1.")
	flag.IntVar(&args.Timeout, "o", 5, "The timeout for API requests. Defaults to 5 seconds.")
	
	// Custom parsing for verbosity
	var vFlag, vvFlag bool
	flag.BoolVar(&vFlag, "v", false, "Enable verbosity level 1.")
	flag.BoolVar(&vvFlag, "vv", false, "Enable verbosity level 2.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] host\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Handle verbosity
	if vvFlag {
		args.Verbosity = 2
	} else if vFlag {
		args.Verbosity = 1
	} else {
		args.Verbosity = 0
	}

	// Positional argument for host
	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "Error: host argument is required.")
		flag.Usage()
		os.Exit(1)
	} else {
		args.Host = flag.Args()[0]
	}

	return args
}