package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

type TestGroup struct {
	Tests []LogTest `json:"Tests"`
}

type LogTest struct {
	Version         string            `json:"Version"`
	RuleID          string            `json:"RuleID"`
	RuleLevel       string            `json:"RuleLevel"`
	RuleDescription string            `json:"RuleDescription"`
	LogFilePath     string            `json:"LogFilePath"`
	Format          string            `json:"Format"`
	Decoder         map[string]string `json:"Decoder"`
	Predecoder      map[string]string `json:"Predecoder"`
	TestDescription string            `json:"TestDescription"`
}

func NewLogTest(Version string, RuleID string, RuleLevel string, RuleDescription string, LogFilePath string, Format string, Decoder map[string]string, Predecoder map[string]string, TestDescription string) (*LogTest, bool, []string, []string) {
	lt := new(LogTest)

	validTest := true // Want to print out all invalid parts of the test
	var errors []string
	var warnings []string

	// Version
	valid, err, warn := isValidVersion(Version)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.Version = Version

	// Rule ID
	valid, err, warn = isValidRuleID(RuleID)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.RuleID = RuleID

	// Rule Level
	valid, err, warn = isValidRuleLevel(RuleLevel)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.RuleLevel = RuleLevel

	// Rule Description
	valid, err, warn = isValidRuleDescription(RuleDescription)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.RuleDescription = RuleDescription

	// Log File Path
	valid, err, warn = isValidLogFilePath(LogFilePath)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.LogFilePath = LogFilePath

	// Format
	valid, err, warn = isValidFormat(Format)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.Format = Format

	// Decoder
	valid, err, warn = isValidDecoder(Decoder)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.Decoder = Decoder

	// Predecoder
	valid, err, warn = isValidPredecoder(Predecoder)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.Predecoder = Predecoder

	// Test Description
	valid, err, warn = isValidTestDescription(TestDescription)
	errors = append(errors, err...)
	warnings = append(warnings, warn...)
	if !valid {
		validTest = false
	}
	lt.TestDescription = TestDescription

	return lt, validTest, errors, warnings
}

func isValidVersion(Version string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	validVersions := map[string]struct{}{
		"0.1": {},
	}

	// Empty Version is not an error but should generally be avoided
	// as it will be assumed to be the latest version
	if Version == "" {
		warnings = append(warnings, "Version is empty using latest test version")
		return true, errors, warnings
	}

	// Check if the version is valid
	_, exists := validVersions[Version]

	if !exists {
		errors = append(errors, fmt.Sprintf("Invalid test version: %s", Version))
		return false, errors, warnings
	}

	return true, errors, warnings
}

// This comes from the Wazuh documentation
// Must be between 0 and 999999
// See: https://documentation.wazuh.com/current/user-manual/ruleset/ruleset-xml-syntax/rules.html#rules-rule
func isValidRuleID(RuleID string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	// Convert RuleID to int
	RuleIDInt, err := strconv.Atoi(RuleID)

	if err != nil {
		errors = append(errors, "Rule ID is not an integer")
		return false, errors, warnings
	}

	if RuleIDInt < 0 {
		errors = append(errors, "Rule ID cannot be less than 0")
		return false, errors, warnings
	}

	if RuleIDInt > 999999 {
		errors = append(errors, "Rule ID cannot be greater than 999999")
		return false, errors, warnings
	}

	return true, errors, warnings
}

// This comes from the Wazuh documentation
// Must be between 0 and 16
// See: https://documentation.wazuh.com/current/user-manual/ruleset/ruleset-xml-syntax/rules.html#rules-rule
func isValidRuleLevel(level string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	// Convert RuleLevel to int
	levelInt, err := strconv.Atoi(level)

	if err != nil {
		errors = append(errors, "Rule level is not an integer")
		return false, errors, warnings
	}

	if levelInt < 0 {
		// Cant be less than 0
		errors = append(errors, "Rule level cannot be less than 0")
		return false, errors, warnings
	}

	if levelInt > 16 {
		// Cant be greater than 16
		errors = append(errors, "Rule level cannot be greater than 16")
		return false, errors, warnings
	}

	return true, errors, warnings
}

// Check that the rule description is not empty
// Generally, the rule description should have some content
func isValidRuleDescription(RuleDescription string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	if RuleDescription == "" {
		warnings = append(warnings, "Rule description is empty")
		return false, errors, warnings
	}

	return true, errors, warnings
}

// Check that the log file path is not empty, file exists,
// is readable, is not emtpy, and has only one line
func isValidLogFilePath(LogFilePath string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	// Check if the log file path is empty
	if LogFilePath == "" || len(LogFilePath) == 0 {
		errors = append(errors, "Log file path is empty")
		return false, errors, warnings
	}

	// Check if the log file exists
	_, err := os.Stat(LogFilePath)
	if err != nil {
		errors = append(errors, "Log file does not exist")
		return false, errors, warnings
	}

	// Check if the log file is readable
	file, err := os.Open(LogFilePath)
	if err != nil {
		errors = append(errors, "Log file is not readable")
		return false, errors, warnings
	}

	// Check if the log file is empty
	stat, err := file.Stat()

	if err != nil {
		errors = append(errors, "Error reading log file size")
		return false, errors, warnings
	}

	if stat.Size() == 0 {
		errors = append(errors, "Log file is empty")
		return false, errors, warnings
	}

	// Check if the log file has only one line
	lineCount, err := fileHasOneLine(file)

	if err != nil {
		errors = append(errors, "Error reading log file lines")
		return false, errors, warnings
	}

	if lineCount < 1 {
		errors = append(errors, "Log file should have at least one line")
		return false, errors, warnings
	}
	if lineCount > 1 {
		errors = append(errors, "Log file has more than one line")
		return false, errors, warnings
	}

	return true, errors, warnings
}

// Reads a file and counts the number of lines in it
// up to a maximum of 2 lines. If the file has more than
// 1 line, it returns 2. If the file has 1 line, it returns 1.
// If the file has no lines, it returns 0.
func fileHasOneLine(r io.Reader) (int, error) {
	buf := make([]byte, 1024)
	count := 0
	lineSep := []byte{'\n'}
	hasData := false

	for {
		c, err := r.Read(buf)
		if c > 0 {
			hasData = true
			count += bytes.Count(buf[:c], lineSep)

			// If we've counted 2 or more lines, return 2
			if count >= 2 {
				return 2, nil
			}
		}

		if err != nil {
			if err == io.EOF {
				// If the file has data but no newline characters, it's considered one line
				if hasData && count == 0 {
					return 1, nil
				}
				return count, nil
			}
			return count, err
		}
	}
}

func isValidFormat(format string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	// See below `log_format` in the API reference:
	// https://documentation.wazuh.com/current/user-manual/api/reference.html#operation/api.controllers.logtest_controller.run_logtest_tool
	validLogTypes := []string{
		"syslog",
		"json",
		"snort-full",
		"squid",
		"eventlog",
		"eventchannel",
		"audit",
		"mysql_log",
		"postgresql_log",
		"nmapg",
		"iis",
		"command",
		"full_command",
		"djb-multilog",
		"multi-line",
	}

	// Empty Format is not an error but should generally be avoided
	if format == "" {
		warnings = append(warnings, "Format is empty")
		return false, errors, warnings
	}

	// Check if the format is valid
	valid := false
	for _, logType := range validLogTypes {
		if format == logType {
			valid = true
			break
		}
	}

	if !valid {
		errors = append(errors, fmt.Sprintf("Log format: %s is not valid", format))
		return false, errors, warnings
	}

	return true, errors, warnings
}

// Checks if any of the decoder values are empty
func isValidDecoder(decoder map[string]string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	// Iterate over map and check if any of the values are empty
	// This is generally a mistake but will not cause a test to
	// break Wazuh Log Test
	for key, value := range decoder {
		if value == "" {
			warnings = append(warnings, "Decoder value for key "+key+" is empty")
			return false, errors, warnings
		}
	}

	return true, errors, warnings
}

// Checks if any of the predecoder values are empty
func isValidPredecoder(predecoder map[string]string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	// Iterate over map and check if any of the values are empty
	// This is generally a mistake but will not cause a test to
	// break Wazuh Log Test
	for key, value := range predecoder {
		if value == "" {
			warnings = append(warnings, "Predecoder value for key "+key+" is empty")
			return false, errors, warnings
		}
	}

	return true, errors, warnings
}

// Checks if test description is empty
func isValidTestDescription(TestDescription string) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}

	// Will not cause a test to break Wazuh Log Test
	// but should generally be avoided
	if TestDescription == "" {
		warnings = append(warnings, "Test description is empty")
		return false, errors, warnings
	}

	return true, errors, warnings
}

func (lt *LogTest) getRuleID() string {
	return lt.RuleID
}

func (lt *LogTest) getTestDescription() string {
	return lt.TestDescription
}

func (lt *LogTest) getLogFilePath() string {
	return lt.LogFilePath
}

func (lt *LogTest) getFormat() string {
	return lt.Format
}

func (lt *LogTest) getRuleDescription() string {
	return lt.RuleDescription
}

func (lt *LogTest) getRuleLevel() string {
	return lt.RuleLevel
}

func (lt *LogTest) getDecoder() map[string]string {
	return lt.Decoder
}

func (lt *LogTest) getPredecoder() map[string]string {
	return lt.Predecoder
}
