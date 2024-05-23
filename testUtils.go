package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// Here is an example structure of:
// rootTestDir/
//
//	test_3.json
//	test_3.log
//	- testGroup1/
//	  test_1.json
//	  test_1.log
//	  test_2.json
//	  test_2.log
//
// Each directory groups tests together. Each file that matches the format
// test_*.json is a test definition/declaration. Each test definition file
// is a list of json logTest objects. Each test will have a corresponding
// log file that is the raw log data that will be sent to Wazuh for processing.
// This can be named anything except test_*.json as it's name defined in the
// test definition file.
//
// Groups can be nested to any depth. The test runner will recursively search
// for test definition files and log files in the root directory and all
// subdirectories.
func runTestGroup(rootTestDir string, numThreads int, verbosity int, timeout int) ([]LogTest, error) {

	// Check if rootTestDir exists
	exists, err := fileExists(rootTestDir)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("root test directory does not exist")
	}

	// Check if rootTestDir is a directory
	isDir, err := isDir(rootTestDir)
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, errors.New("root test directory is not a directory")
	}

	// List our current directory
	files, err := os.ReadDir(rootTestDir)
	if err != nil {
		return nil, err
	}

	// Sort the files into test definitions, subdirectories, and other files
	testDefs, subdirectories, otherFiles, err := sortDirContent(files)

	if err != nil {
		return nil, err
	}

	// Recurse into subdirectories to evaluate tests
	// test tree from bottom up
	for _, subdirectory := range subdirectories {
		path := filepath.Join(rootTestDir, subdirectory.Name())
		_, err := runTestGroup(path, numThreads, verbosity, timeout)
		if err != nil {
			return nil, err
		}
	}

	// Load all test definitions
	var logTests []LogTest
	for _, testDef := range testDefs {
		path := filepath.Join(rootTestDir, testDef.Name())
		tests, err := loadTestDef(path, verbosity)
		if err != nil {
			panic(err)
		}
		logTests = append(logTests, tests...)
	}

	if len(logTests) > 0 && verbosity > 0 {
		PrintBoldWhite("Loaded " + strconv.Itoa(len(logTests)) + " tests from " + rootTestDir)
	}

	// There should be a one to one mapping of
	// LogTest objects to log files. This means
	// that the number of other files should be
	// greater than or equal to the number of LogTest
	// objects for the current directory.
	//
	// Warn the users upfront so they are aware
	// when interpreting the results.
	if len(otherFiles) < len(logTests) {
		diff := len(logTests) - len(otherFiles)
		PrintYellow("WARNING: " + rootTestDir + " has " + strconv.Itoa(diff) + " tests with no corresponding log files and will be skipped...")
	}

	// Create progres bar for visual feedback
	bar := progressbar.NewOptions(len(logTests), progressbar.OptionSetDescription("Running: "+rootTestDir), progressbar.OptionShowCount())

	// Create a buffered channel to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, numThreads)
	var wg sync.WaitGroup

	// Save failed tests for reporting
	var failedTests []LogTest
	var failTestErrors [][]string

	// Run tests concurrently with a max of user-defined number of threads
	for _, logTest := range logTests {
		wg.Add(1)
		semaphore <- struct{}{} // acquire a slot
		go func(logTest LogTest) {
			defer wg.Done()
			defer func() { <-semaphore }() // release the slot
			passed, testErrors := runTest(logTest, rootTestDir)

			if !passed {
				failedTests = append(failedTests, logTest)
				failTestErrors = append(failTestErrors, testErrors)
			}

			bar.Add(1)
		}(logTest)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Printf("\n\n")

	return logTests, nil
}

// This function will run a single test and return back the pass/fail
// and any errors that occurred during the test.
func runTest(logTest LogTest, workingDir string) (bool, []string) {

	// TODO: Implement the test runner

	return true, nil
}

func loadTestDef(path string, verbosity int) ([]LogTest, error) {
	// Check file extension is .json
	if filepath.Ext(path) != ".json" {
		return nil, errors.New("file is not a JSON file")
	}

	// Check if path exists
	exists, err := fileExists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("file does not exist")
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse JSON list of LogTest objects
	var testGroup TestGroup
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&testGroup)
	if err != nil {
		return nil, err
	}

	var invalidTestCount int
	var logTests []LogTest
	for i, raw := range testGroup.Tests {
		logPath := filepath.Join(filepath.Dir(path), raw.LogFilePath)
		logTest, valid, loadErrors, loadWarnings := NewLogTest(raw.RuleID, raw.RuleLevel, raw.RuleDescription, logPath, raw.Format, raw.Decoder, raw.Predecoder, raw.TestDescription)
		if !valid {
			// Print warnings or handle invalid tests as needed
			if logTest.getRuleID() == "" {
				PrintWhite("[FAILED LOAD] " + path + ": Test #" + strconv.Itoa(i+1))
			} else {
				PrintWhite("[FAILED LOAD] Test: (RuleID: " + logTest.getRuleID() + ") " + logTest.getTestDescription())
			}

			invalidTestCount++

			if verbosity < 1 {
				continue
			}

			// Print Errors for: -v (1), -vv (2)
			// Tab over to show that these are errors
			// corresponding to the test above
			if verbosity > 0 && len(loadErrors) > 0 {
				for _, e := range loadErrors {
					PrintRed("+ " + e)
				}
			}

			// Print Warnings for: -vv (2)
			if verbosity > 1 && len(loadWarnings) > 0 {
				for _, w := range loadWarnings {
					PrintYellow("+ " + w)
				}
			}

			fmt.Printf("\n")

		}

		logTests = append(logTests, *logTest)
	}

	return logTests, nil
}

// Load all test definitions from the current directory
func sortDirContent(files []os.DirEntry) ([]os.DirEntry, []os.DirEntry, []os.DirEntry, error) {
	var testDefs []os.DirEntry
	var subdirectories []os.DirEntry
	var otherFiles []os.DirEntry

	// Compile the regex pattern
	pattern := regexp.MustCompile(`^test_.*\.json$`)

	for _, file := range files {
		if file.IsDir() {
			subdirectories = append(subdirectories, file)
			continue
		}

		if pattern.MatchString(file.Name()) {
			testDefs = append(testDefs, file)
			continue
		}

		// If the file is not a directory or a test definition, then it is
		// some other file; Likely a log file
		otherFiles = append(otherFiles, file)
	}

	return testDefs, subdirectories, otherFiles, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func isDir(path string) (bool, error) {

	// First check if path does not exist
	exists, err := fileExists(path)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	// Check if path is a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
