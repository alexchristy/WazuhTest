package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
func runTestGroup(ws WazuhServer, rootTestDir string, numThreads int, verbosity int, timeout int) ([]LogTest, error) {

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
		_, err := runTestGroup(ws, path, numThreads, verbosity, timeout)
		if err != nil {
			return nil, err
		}
	}

	// Load all test definitions
	var logTests []LogTest = []LogTest{}
	var invalidTests int = 0
	for _, testDef := range testDefs {
		path := filepath.Join(rootTestDir, testDef.Name())
		tests, currInvalidTests, err := loadTestDef(path, verbosity)

		// This panic will only occur
		// if the test definition file (.json)
		// has an error.
		if err != nil {
			panic(err)
		}
		logTests = append(logTests, tests...)
		invalidTests += currInvalidTests
	}

	if len(logTests) > 0 && verbosity > 0 {
		totalTests := len(logTests) + invalidTests
		PrintBoldWhite("Sucessfully loaded " + strconv.Itoa(len(logTests)) + "/" + strconv.Itoa(totalTests) + " tests from " + rootTestDir)
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
	var failedTest int = 0
	var errors [][]string
	var warnings [][]string

	// Run tests concurrently with a max of user-defined number of threads
	for _, logTest := range logTests {
		wg.Add(1)
		semaphore <- struct{}{} // acquire a slot
		go func(logTest LogTest) {
			defer wg.Done()
			defer func() { <-semaphore }() // release the slot
			passed, testErrors, testWarnings := runTest(ws, logTest)

			if !passed {
				failedTest++
			}

			errors = append(errors, testErrors)
			warnings = append(warnings, testWarnings)

			bar.Add(1)
		}(logTest)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Printf("\n")

	// len(logTests) == len(errors) == len(warnings)
	for i, test := range logTests {
		failedTest := false

		if len(errors[i]) > 0 {
			failedTest = true
			PrintRed("[FAILED] Test: (RuleID: " + test.getRuleID() + ") " + test.getTestDescription())
			for _, e := range errors[i] {
				PrintRed("+ " + e)
			}
		}
		if verbosity > 1 && len(warnings[i]) > 0 {
			// Only print warnings header if there were no errors
			if !failedTest {
				PrintYellow("[WARNING] Test: (RuleID: " + test.getRuleID() + ") " + test.getTestDescription())
			}
			for _, w := range warnings[i] {
				PrintYellow("+ " + w)
			}
		}
	}

	return logTests, nil
}

// This function will run a single test and return back the pass/fail
// and any errors that occurred during the test.
func runTest(ws WazuhServer, logTest LogTest) (bool, []string, []string) {

	var errors []string
	var warnings []string

	// Load the log file
	logData, err := os.ReadFile(logTest.getLogFilePath())
	if err != nil {
		errors = append(errors, "Error opening log file: "+err.Error())
		return false, errors, warnings
	}

	// Create headers for request
	logTestHeaders := map[string]interface{}{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + ws.getToken(),
	}

	// Create data to send with request
	logTestData := map[string]interface{}{
		"event":      string(logData),
		"log_format": logTest.getFormat(),
		"location":   "WazuhTestRunner",
	}
	jsonData, err := json.Marshal(logTestData)
	if err != nil {
		errors = append(errors, "Error marshalling log data: "+err.Error())
		return false, errors, warnings
	}

	// Build request to send logTestData
	req, err := http.NewRequest("PUT", ws.getLogTestUrl(), bytes.NewBuffer(jsonData))
	if err != nil {
		errors = append(errors, "Error creating request: "+err.Error())
		return false, errors, warnings
	}

	// Send request
	result, err := ws.sendRequest(req, logTestHeaders)
	if err != nil {
		errors = append(errors, "Error sending request: "+err.Error())
		return false, errors, warnings
	}

	// Convert result map to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Error marshalling Wazuh server result map to JSON: %v", err)
	}

	// Unmarshal JSON bytes into the Response struct
	var response Response
	err = json.Unmarshal(jsonBytes, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling Wazuh server response JSON to Response struct: %v", err)
	}

	// Validate the response
	passed, resErrors, resWarnings := validateLogTestResponse(logTest, response)
	if !passed {
		errors = append(errors, resErrors...)
		warnings = append(warnings, resWarnings...)
	}

	return passed, errors, warnings
}

// This function will compare the expected response
// with the actual response from the Wazuh server.
func validateLogTestResponse(logTest LogTest, response Response) (bool, []string, []string) {
	// Note: All the errors and warnings
	// that we recieve from the validation
	// functions are related to the value
	// that is returned from the Wazuh server.
	//
	// At this point the LogTests have already
	// been validated for correctness.
	//
	// We will return early for all errors
	// to prevent giving back more errors
	// than necessary.
	var errors []string
	var warnings []string
	var passed bool = true

	// ======( RuleID Validation )====== //
	passed, ruleIDErrors, ruleIDWarnings := validateRuleID(logTest.getRuleID(), response.Data.Output.Rule.ID)
	if !passed {
		errors = append(errors, ruleIDErrors...)
		warnings = append(warnings, ruleIDWarnings...)
		return passed, errors, warnings
	}

	// ======( RuleLevel Validation )====== //
	expectedRuleLevel, err := strconv.Atoi(logTest.getRuleLevel())
	if err != nil {
		errors = append(errors, "Error converting returned RuleLevel to int: "+err.Error())
		return false, errors, warnings
	}

	passed, ruleLevelErrors, ruleLevelWarnings := validateRuleLevel(expectedRuleLevel, response.Data.Output.Rule.Level)
	if !passed {
		errors = append(errors, ruleLevelErrors...)
		warnings = append(warnings, ruleLevelWarnings...)
		return passed, errors, warnings
	}

	// ======( RuleDescription Validation )====== //
	passed, ruleDescriptionErrors, ruleDescriptionWarnings := validateRuleDescription(logTest.getRuleDescription(), response.Data.Output.Rule.Description)
	if !passed {
		errors = append(errors, ruleDescriptionErrors...)
		warnings = append(warnings, ruleDescriptionWarnings...)
		return passed, errors, warnings
	}

	// ======( Predecoder Validation )====== //
	passed, predecoderErrors, predecoderWarnings := validateDecoder(logTest.getPredecoder(), response.Data.Output.Predecoder, "Predecoder")
	if !passed {
		errors = append(errors, predecoderErrors...)
		warnings = append(warnings, predecoderWarnings...)
		return passed, errors, warnings
	}

	// ======( Decoder Validation )====== //
	passed, decoderErrors, decoderWarnings := validateDecoder(logTest.getDecoder(), response.Data.Output.Decoder, "Decoder")
	if !passed {
		errors = append(errors, decoderErrors...)
		warnings = append(warnings, decoderWarnings...)
		return passed, errors, warnings
	}

	return passed, errors, warnings
}

// This function will check for messages in
func extractLogTestMessages(result map[string]interface{}) (bool, []string, []string) {
	// Based on what I can find from the documentation
	// and personal experience, there appear to be
	// only INFO and WARNING messages.
	var warnings []string
	var info []string
	var haveMessages bool = false

	// Check if the messages key exists
	messages, ok := result["data"].(map[string]interface{})["messages"]
	if !ok {
		return false, nil, nil
	}

	// Check if the messages key is empty
	if len(messages.([]interface{})) == 0 {
		return false, nil, nil
	}

	// Extract the messages
	// regexp.MustCompile()
	// for _, message := range messages.([]interface{}) {

	return haveMessages, warnings, info
}

// This function will validate the RuleID returned by the Wazuh server
func validateRuleID(expected string, got string) (bool, []string, []string) {
	var errors []string
	var warnings []string

	// Note: we are assuming that
	// expected is correct and would
	// pass a isValidRuleID check.
	//
	// We check first if the recieved
	// RuleID is invalid to provide
	// better feedback to the user.
	//
	// We also return early if the RuleID
	// is invalid to prevent giving back
	// more errors than necessary.
	if got == "" {
		errors = append(errors, "RuleID is empty")
		return false, errors, warnings

	}

	// Check if the returned RuleID is valid
	valid, valErrors, valWarnigns := isValidRuleID(got)
	errors = append(errors, valErrors...)
	warnings = append(warnings, valWarnigns...)
	if !valid {
		return false, errors, warnings
	}

	// Check if the expected RuleID matches the
	// returned RuleID
	if expected != got {
		errors = append(errors, "Expected RuleID: "+expected+" Got RuleID: "+got)
		return false, errors, warnings
	}

	return true, errors, warnings
}

// This function will validate the RuleLevel returned by the Wazuh server
func validateRuleLevel(expected int, got int) (bool, []string, []string) {
	var errors []string
	var warnings []string

	// Check if the RuleLevel is valid
	valid, valErrors, valWarnings := isValidRuleLevel(strconv.Itoa(got))
	errors = append(errors, valErrors...)
	warnings = append(warnings, valWarnings...)
	if !valid {
		return false, errors, warnings
	}

	// Check if the expected RuleLevel matches the
	// returned RuleLevel
	if expected != got {
		errors = append(errors, "Expected RuleLevel: "+strconv.Itoa(expected)+" Got RuleLevel: "+strconv.Itoa(got))
		return false, errors, warnings
	}

	return true, errors, warnings
}

// This function will validate the RuleDescription returned by the Wazuh server
func validateRuleDescription(expected string, got string) (bool, []string, []string) {
	var errors []string
	var warnings []string

	// =====( RuleDescription Validation )=====
	if got == "" {
		errors = append(errors, "RuleDescription is empty")
		return false, errors, warnings
	}

	// Check if the RuleDescription is valid
	valid, valErrors, valWarnings := isValidRuleDescription(got)
	errors = append(errors, valErrors...)
	warnings = append(warnings, valWarnings...)
	if !valid {
		return false, errors, warnings
	}

	// Check if the expected RuleDescription matches the
	// returned RuleDescription
	if expected != got {
		errors = append(errors, "Expected RuleDescription: "+expected+" Got RuleDescription: "+got)
		return false, errors, warnings
	}

	return true, errors, warnings
}

func validateDecoder(expected map[string]string, got map[string]string, decoderType string) (bool, []string, []string) {
	var errors []string
	var warnings []string
	var passed bool = true

	// Check if the expected Decoder is empty
	// this means that the test does not care
	// about the Decoder output.
	if len(expected) == 0 {
		return true, errors, warnings
	}

	for key, val := range expected {
		// Check if the key exists in the returned Decoder
		_, ok := got[key]
		if !ok {
			passed = false
			errors = append(errors, "Expected key: "+key+" not found in returned "+decoderType)
			continue
		}

		// Check if the value of the key matches the expected value
		if val != got[key] {
			passed = false
			errors = append(errors, "Expected value: "+val+" for key: "+key+" in returned "+decoderType+" Got value: "+got[key])
			continue
		}
	}

	return passed, errors, warnings
}

func loadTestDef(path string, verbosity int) ([]LogTest, int, error) {
	var invalidTestCount int = 0

	// Check file extension is .json
	if filepath.Ext(path) != ".json" {
		return nil, -1, errors.New("file is not a JSON file")
	}

	// Check if path exists
	exists, err := fileExists(path)
	if err != nil {
		return nil, -1, err
	}
	if !exists {
		return nil, 1, errors.New("file does not exist")
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, -1, err
	}
	defer file.Close()

	// Parse JSON list of LogTest objects
	var testGroup TestGroup
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&testGroup)
	if err != nil {
		return nil, -1, err
	}

	var logTests []LogTest
	for i, raw := range testGroup.Tests {
		logPath := filepath.Join(filepath.Dir(path), raw.LogFilePath)
		logTest, valid, loadErrors, loadWarnings := NewLogTest(raw.Version, raw.RuleID, raw.RuleLevel, raw.RuleDescription, logPath, raw.Format, raw.Decoder, raw.Predecoder, raw.TestDescription)
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

		// Do not append invalid tests
		if valid {
			logTests = append(logTests, *logTest)
		}
	}

	return logTests, invalidTestCount, nil
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
