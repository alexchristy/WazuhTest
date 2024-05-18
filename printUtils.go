package main

import (
	"fmt"
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
