package main

import (
	"fmt"
	"strconv"
)

func PrintRed(text string) {
	fmt.Println("\033[91m" + text + "\033[0m")
}

func PrintGreen(text string) {
	fmt.Println("\033[92m" + text + "\033[0m")
}

func PrintYellow(text string) {
	fmt.Println("\033[93m" + text + "\033[0m")
}

func PrintWhite(text string) {
	fmt.Println("\033[97m" + text + "\033[0m")
}

func PrintBoldWhite(text string) {
	fmt.Println("\033[97m\033[1m" + text + "\033[0m")
}

func printSummary(numTests int, numFailedTests int, numWarnTests int) error {

	PrintBoldWhite("Test Summary:")
	PrintBoldWhite("=============\n")

	fmt.Printf("Total: %d\n", numTests)

	if numFailedTests > 0 {
		PrintRed("Failed: " + strconv.Itoa(numFailedTests))
	}

	if numWarnTests > 0 {
		PrintYellow("Warned: " + strconv.Itoa(numWarnTests))
	}

	fmt.Printf("\n")

	if numFailedTests <= 0 {
		PrintGreen("All tests passed.")
	}

	return nil
}
