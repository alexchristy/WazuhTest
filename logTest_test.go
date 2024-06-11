package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Test_isValidVersion(t *testing.T) {
	type args struct {
		Version string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{
		// Valid Versions
		{name: "Valid version 0.1", args: args{"0.1"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid but warn for blank version", args: args{""}, want: true, want1: []string{}, want2: []string{"Version is empty using latest test version"}},

		// Invalid Versions
		{name: "Invalid negative version", args: args{"-0.1"}, want: false, want1: []string{"Invalid test version: -0.1"}, want2: []string{}},
		{name: "Invalid non-numeric version", args: args{"hello"}, want: false, want1: []string{"Invalid test version: hello"}, want2: []string{}},
		{name: "Invalid non-existent version", args: args{"100.0.1"}, want: false, want1: []string{"Invalid test version: 100.0.1"}, want2: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Logf(fmt.Sprintf("%+v", tt))
			got, got1, got2 := isValidVersion(tt.args.Version)
			if got != tt.want {
				t.Errorf("isValidVersion() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidVersion() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidVersion() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Fuzz_isValidVersion(f *testing.F) {
	// Adding a string instead of an integer
	f.Add("10167")
	f.Fuzz(func(t *testing.T, v string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occurred: %v", r)
			}
		}()

		isValidVersion(v)
	})
}

func Benchmark_isValidVersion_validVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidVersion("0.1")
	}
}

func Benchmark_isValidVersion_invalidVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidVersion(strconv.Itoa(i))
	}
}

func Test_isValidRuleID(t *testing.T) {
	type args struct {
		RuleID string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{
		// Valid ruleID's
		{name: "Valid lowest RuleID 0", args: args{"0"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid highest RuleID 999999", args: args{"999999"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid RuleID 101010", args: args{"101010"}, want: true, want1: []string{}, want2: []string{}},

		// Invalid ruleID's
		{name: "Invalid negative RuleID -15", args: args{"-15"}, want: false, want1: []string{"Invalid rule ID cannot be less than 0"}, want2: []string{}},
		{name: "Invalid too large RuleID 1234567890", args: args{"1234567890"}, want: false, want1: []string{"Invalid rule ID cannot be greater than 999999"}, want2: []string{}},
		{name: "Invalid non-numeric RuleID hello", args: args{"hello"}, want: false, want1: []string{"Invalid rule ID is not an integer"}, want2: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := isValidRuleID(tt.args.RuleID)
			if got != tt.want {
				t.Errorf("isValidRuleID() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidRuleID() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidRuleID() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Fuzz_isValidRuleID(f *testing.F) {
	// Adding a string instead of an integer
	f.Add("10167")
	f.Fuzz(func(t *testing.T, id string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occurred: %v", r)
			}
		}()

		isValidRuleID(id)
	})
}
func Benchmark_isValidRuleID_validRuleID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidRuleID("999")
	}
}

func Benchmark_isValidRuleID_invalidRuleID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidRuleID(strconv.Itoa(i))
	}
}

func Test_isValidRuleLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{
		// Valid rule levels
		{name: "Valid lowest rule level 0", args: args{"0"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid highest rule level 16", args: args{"16"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid middle rule level 8", args: args{"8"}, want: true, want1: []string{}, want2: []string{}},

		// Invalid rule levels
		{name: "Invalid negative rule level -6", args: args{"-6"}, want: false, want1: []string{"Rule level cannot be less than 0"}, want2: []string{}},
		{name: "Invalid too large rule level 600", args: args{"600"}, want: false, want1: []string{"Rule level cannot be greater than 16"}, want2: []string{}},
		{name: "Invalid non-numeric rule level hello", args: args{"hello"}, want: false, want1: []string{"Rule level is not an integer"}, want2: []string{}},
		{name: "Invalid empty rule level", args: args{""}, want: false, want1: []string{"Rule level is not an integer"}, want2: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := isValidRuleLevel(tt.args.level)
			if got != tt.want {
				t.Errorf("isValidRuleLevel() got = %v, want %v", got, tt.want)
			}

			if len(tt.want1) > 0 { // If we are expecting errors
				hasError := false
				for _, err := range got1 { // errors
					err := strings.ToLower(err)
					if strings.Contains(err, "invalid") && strings.Contains(err, "rule level") {
						hasError = true
						break
					}
				}
				if !hasError {
					t.Errorf("Expected to get error for invalid rule level")
				}
			}

			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidRuleLevel() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Fuzz_isValidRuleLevel(f *testing.F) {
	examples := []string{"0", "16", "1000", "-373"}
	for _, ex := range examples {
		f.Add(ex)
	}

	f.Fuzz(func(t *testing.T, l string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occured: %v", r)
			}
		}()

		isValidRuleLevel(l)
	})
}

func Benchmark_isValidRuleLevel_validRuleLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidRuleLevel("15")
	}
}

func Benchmark_isValidRuleLevel_invalidRuleLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidRuleLevel(strconv.Itoa(i))
	}
}

func Test_isValidRuleDescription(t *testing.T) {
	type args struct {
		RuleDescription string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{
		{name: "Valid description", args: args{"This is a description."}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid emoji description", args: args{"🥝🥝🥝🥝"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid empty description with warning", args: args{""}, want: true, want1: []string{}, want2: []string{"Rule description is empty"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := isValidRuleDescription(tt.args.RuleDescription)
			if got != tt.want {
				t.Errorf("isValidRuleDescription() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidRuleDescription() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidRuleDescription() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Fuzz_isValidRuleDescription(f *testing.F) {
	examples := []string{"This is a test description.", "It can have unicode 👍Ư", "More unicode ؿ", "\"\"\\"}
	for _, ex := range examples {
		f.Add(ex)
	}

	f.Fuzz(func(t *testing.T, d string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occured: %v", r)
			}
		}()

		isValidRuleDescription(d)
	})
}

func Benchmark_isValidRuleDescription_validDescription(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidRuleDescription("This is a valid rule description")
	}
}

func Benchmark_isValidRuleDescription_validDescriptionWithWarn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidRuleDescription("")
	}
}

func Test_isValidLogFilePath(t *testing.T) {
	// Create a file for testing purposes
	files := map[string]string{
		"emptyFileName":    "./empty.log",
		"zeroLineFileName": "./zero_line.log",
		"oneLineFileName":  "./one_line.log",
		"twoLineFileName":  "./two_line.log",
	}

	// Delete the test files when complete
	t.Cleanup(func() {
		for _, fileName := range files {
			err := os.Remove(fileName)
			if err != nil {
				t.Fatalf("Error deleting file: %s", fileName)
			}
		}
	})

	// Empty file
	_, err := os.Create(files["emptyFileName"])
	if err != nil {
		t.Fatalf("Failed to create empty test file")
	}

	// Zero line log file
	zeroLineFile, err := os.Create(files["zeroLineFileName"])
	if err != nil {
		t.Fatalf("Failed to create zero line test file")
	}
	// With no newline at the end `cat zero_line.log | wc -l` returns 0
	zeroLineStr := "This is a zero line log"
	_, err = zeroLineFile.WriteString(zeroLineStr)
	if err != nil {
		t.Fatalf("Failed to write zero line log to test file")
	}

	// One line log file
	oneLineFile, err := os.Create(files["oneLineFileName"])
	if err != nil {
		t.Fatalf("Failed to create one line test file")
	}
	oneLineStr := "This is a one line log\n"
	_, err = oneLineFile.WriteString(oneLineStr)
	if err != nil {
		t.Fatalf("Failed to write one line log to test file")
	}

	// Two line log file
	twoLineFile, err := os.Create(files["twoLineFileName"])
	if err != nil {
		t.Fatalf("Failed to create two line test file")
	}
	twoLineStr := "This is a two line log\nThis is the second line of the log\n"
	_, err = twoLineFile.WriteString(twoLineStr)
	if err != nil {
		t.Fatalf("Failed to write two line log to test file")
	}

	type args struct {
		LogFilePath string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{
		// Valid log files
		{name: "Valid single line log file", args: args{files["oneLineFileName"]}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid single line log file no newline termination", args: args{files["zeroLineFileName"]}, want: true, want1: []string{}, want2: []string{}},

		// Invalid log files
		{name: "Invalid empty log file", args: args{files["emptyFileName"]}, want: false, want1: []string{"Log file is empty"}, want2: []string{}},
		{name: "Invalid empty log file path", args: args{""}, want: false, want1: []string{"Log file path is empty"}, want2: []string{}},
		{name: "Invalid non-existent log file path", args: args{"./i-dont-exist.log"}, want: false, want1: []string{"Log file does not exist"}, want2: []string{}},
		{name: "Invalid multi-line log file", args: args{files["twoLineFileName"]}, want: false, want1: []string{"Log file should only have one line"}, want2: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := isValidLogFilePath(tt.args.LogFilePath)
			if got != tt.want {
				t.Errorf("isValidLogFilePath() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidLogFilePath() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidLogFilePath() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Fuzz_isValidLogFilePath(f *testing.F) {
	examples := []string{"./log.log", "./wazuh.log", "./file", "/tmp/full/file/path"}
	for _, ex := range examples {
		f.Add(ex)
	}

	f.Fuzz(func(t *testing.T, f string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occured: %v", r)
			}
		}()

		isValidLogFilePath(f)
	})
}

func Benchmark_isValidLogFilePath_invalidLogFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidLogFilePath(strconv.Itoa(i))
	}
}

func Test_fileHasOneLine(t *testing.T) {
	// Create a file for testing purposes
	files := map[string]string{
		"emptyFileName":    "./empty.log",
		"zeroLineFileName": "./zero_line.log",
		"oneLineFileName":  "./one_line.log",
		"twoLineFileName":  "./two_line.log",
	}

	// Delete the test files when complete
	t.Cleanup(func() {
		for _, fileName := range files {
			err := os.Remove(fileName)
			if err != nil {
				t.Fatalf("Error deleting file: %s", fileName)
			}
		}
	})

	// Empty file
	emptyFile, err := os.Create(files["emptyFileName"])
	if err != nil {
		t.Fatalf("Failed to create empty test file")
	}

	// Zero line log file
	zeroLineFile, err := os.Create(files["zeroLineFileName"])
	if err != nil {
		t.Fatalf("Failed to create zero line test file")
	}
	// With no newline at the end `cat zero_line.log | wc -l` returns 0
	zeroLineStr := "This is a zero line log"
	_, err = zeroLineFile.WriteString(zeroLineStr)
	if err != nil {
		t.Fatalf("Failed to write zero line log to test file")
	}
	_, err = zeroLineFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("Failed to seek zero line test file")
	}

	// One line log file
	oneLineFile, err := os.Create(files["oneLineFileName"])
	if err != nil {
		t.Fatalf("Failed to create one line test file")
	}
	oneLineStr := "This is a one line log\n"
	_, err = oneLineFile.WriteString(oneLineStr)
	if err != nil {
		t.Fatalf("Failed to write one line log to test file")
	}
	// See back to the beginning of the file for the test to work properly
	_, err = oneLineFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("Failed to seek one line test file")
	}

	// Two line log file
	twoLineFile, err := os.Create(files["twoLineFileName"])
	if err != nil {
		t.Fatalf("Failed to create two line test file")
	}
	twoLineStr := "This is a two line log\nThis is the second line of the log\n"
	_, err = twoLineFile.WriteString(twoLineStr)
	if err != nil {
		t.Fatalf("Failed to write two line log to test file")
	}
	_, err = twoLineFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("Failed to seek two line test file")
	}

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// Valid single line files
		{name: "Valid single line log file", args: args{oneLineFile}, want: true, wantErr: false},
		{name: "Valid single line log file with no newline terminator", args: args{zeroLineFile}, want: true, wantErr: false},

		// Valid non-single line files
		{name: "Empty file", args: args{emptyFile}, want: false, wantErr: false},
		{name: "Two line file", args: args{twoLineFile}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := fileHasOneLine(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileHasOneLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileHasOneLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidFormat(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{
		// Valid formats
		{name: "Valid format syslog", args: args{"syslog"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid format json", args: args{"json"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid format audit", args: args{"audit"}, want: true, want1: []string{}, want2: []string{}},

		// Invalid formats
		{name: "Invalid format NON_VALID", args: args{"NON_VALID"}, want: false, want1: []string{"Log format: NON_VALID is not valid"}, want2: []string{}},
		{name: "Invalid format unicode 🥝🥝🥝🥝", args: args{"🥝🥝🥝🥝"}, want: false, want1: []string{"Log format: 🥝🥝🥝🥝 is not valid"}, want2: []string{}},
		{name: "Invalid empty format", args: args{""}, want: false, want1: []string{"Format is empty"}, want2: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := isValidFormat(tt.args.format)
			if got != tt.want {
				t.Errorf("isValidFormat() got = %v, want %v", got, tt.want)
			}

			if len(tt.want1) > 0 { // If we are expecting errors
				hasError := false
				for _, err := range got1 { // errors
					err := strings.ToLower(err)
					if strings.Contains(err, "invalid") && strings.Contains(err, "format") {
						hasError = true
						break
					}
				}
				if !hasError {
					t.Errorf("Expected to get error for invalid rule level")
				}
			}

			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidFormat() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Fuzz_isValidFormat(f *testing.F) {
	examples := []string{"syslog", "audit", "multi-log", "bleh*bleh"}
	for _, ex := range examples {
		f.Add(ex)
	}

	f.Fuzz(func(t *testing.T, f string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occured: %v", r)
			}
		}()

		isValidFormat(f)
	})
}

func Benchmark_isValidFormat_invalidFormats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidFormat(strconv.Itoa(i))
	}
}

func Test_isValidDecoder(t *testing.T) {
	type args struct {
		decoder map[string]string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{
		// Valid decoders
		{name: "Valid single element decoder", args: args{map[string]string{"key": "value"}}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid single element unicode decoder", args: args{map[string]string{"🥝🥝": "🥝🥝"}}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid empty decoder with warning", args: args{map[string]string{}}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid multi element decoder", args: args{map[string]string{"key1": "value1", "key2": "value2"}}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid single key empty value decoder with warning", args: args{map[string]string{"emptyKey": ""}}, want: true, want1: []string{}, want2: []string{"Decoder value for key emptyKey is empty"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := isValidDecoder(tt.args.decoder)
			if got != tt.want {
				t.Errorf("isValidDecoder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidDecoder() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidDecoder() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Benchmark_isValidDecoder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		num := strconv.Itoa(i)
		isValidDecoder(map[string]string{num: num})
	}
}

func Test_isValidPredecoder(t *testing.T) {
	type args struct {
		predecoder map[string]string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{ // Valid predecoders
		{name: "Valid single element predecoder", args: args{map[string]string{"key": "value"}}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid single element unicode predecoder", args: args{map[string]string{"🥝🥝": "🥝🥝"}}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid empty predecoder with warning", args: args{map[string]string{}}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid multi element predecoder", args: args{map[string]string{"key1": "value1", "key2": "value2"}}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid single key empty value predecoder with warning", args: args{map[string]string{"emptyKey": ""}}, want: true, want1: []string{}, want2: []string{"Predecoder value for key emptyKey is empty"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := isValidPredecoder(tt.args.predecoder)
			if got != tt.want {
				t.Errorf("isValidPredecoder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidPredecoder() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidPredecoder() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Benchmark_isValidPredecoder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		num := strconv.Itoa(i)
		isValidPredecoder(map[string]string{num: num})
	}
}

func Test_isValidTestDescription(t *testing.T) {
	type args struct {
		TestDescription string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []string
		want2 []string
	}{
		// Valid test descriptions
		{name: "Valid test description", args: args{"This is a valid test description"}, want: true, want1: []string{}, want2: []string{}},
		{name: "Valid empty test description with warning", args: args{""}, want: true, want1: []string{}, want2: []string{"Test description is empty"}},
		{name: "Valid unicode test description", args: args{"🥝🥝🥝🥝"}, want: true, want1: []string{}, want2: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := isValidTestDescription(tt.args.TestDescription)
			if got != tt.want {
				t.Errorf("isValidTestDescription() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidTestDescription() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidTestDescription() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Fuzz_isValidTestDescription(f *testing.F) {
	examples := []string{"syslog", "audit", "multi-log", "bleh*bleh", "This is a description🥝"}
	for _, ex := range examples {
		f.Add(ex)
	}

	f.Fuzz(func(t *testing.T, td string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic occured: %v", r)
			}
		}()

		isValidTestDescription(td)
	})
}

func Benchmark_isValidTestDescription(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidTestDescription(strconv.Itoa(i))
	}
}

func createTestLog(log string) (string, error) {
	logFile, err := os.CreateTemp("", "testLog*")
	if err != nil {
		return "", err
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			panic(err)
		}
	}(logFile)

	_, err = logFile.WriteString(log)
	if err != nil {
		return "", err
	}
	_, err = logFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	return logFile.Name(), nil
}

func Test_NewLogTestValidTests(t *testing.T) {
	type args struct {
		Version         string
		RuleID          string
		RuleLevel       string
		RuleDescription string
		LogFilePath     string
		Format          string
		Decoder         map[string]string
		Predecoder      map[string]string
		TestDescription string
	}

	// Minimal log test
	ltMinFile, err := createTestLog("This is a test log file.")
	if err != nil {
		t.Fatalf("Failed to create minimal LogTest log file")
	}
	ltMinimalArgs := args{
		Version:         "0.1",
		RuleID:          "100",
		RuleLevel:       "8",
		RuleDescription: "This is a description.",
		LogFilePath:     ltMinFile,
		Format:          "syslog",
		Decoder:         map[string]string{},
		Predecoder:      map[string]string{},
		TestDescription: "Minimal LogTest for testing purposes",
	}
	// corr = correct (i.e. valid/will generate)
	corrMinLogTest := new(LogTest)
	corrMinLogTest.Version = "0.1"
	corrMinLogTest.RuleID = "100"
	corrMinLogTest.RuleLevel = "8"
	corrMinLogTest.RuleDescription = "This is a description."
	corrMinLogTest.LogFilePath = ltMinFile
	corrMinLogTest.Format = "syslog"
	corrMinLogTest.Decoder = map[string]string{}
	corrMinLogTest.Predecoder = map[string]string{}
	corrMinLogTest.TestDescription = "Minimal LogTest for testing purposes"

	// All warnings log test
	allWarningArgs := args{
		Version:         "",
		RuleID:          "1000",
		RuleLevel:       "14",
		RuleDescription: "",
		LogFilePath:     ltMinFile,
		Format:          "syslog",
		Decoder: map[string]string{
			"emptyDecoderKey": "",
		},
		Predecoder: map[string]string{
			"emptyPredecoderKey": "",
		},
		TestDescription: "",
	}
	allWarningTest := new(LogTest)
	allWarningTest.Version = ""
	allWarningTest.RuleID = "1000"
	allWarningTest.RuleLevel = "14"
	allWarningTest.RuleDescription = ""
	allWarningTest.LogFilePath = ltMinFile
	allWarningTest.Format = "syslog"
	allWarningTest.Decoder = map[string]string{
		"emptyDecoderKey": "",
	}
	allWarningTest.Predecoder = map[string]string{
		"emptyPredecoderKey": "",
	}

	allExpectedWarns := []string{
		"Version is empty using latest test version",
		"Rule description is empty",
		"Decoder value for key emptyDecoderKey is empty",
		"Predecoder value for key emptyPredecoderKey is empty",
		"Test description is empty",
	}

	tests := []struct {
		name  string
		args  args
		want  *LogTest
		want1 bool
		want2 []string
		want3 []string
	}{
		{name: "Valid minimum valid LogTest", args: ltMinimalArgs, want: corrMinLogTest, want1: true, want2: []string{}, want3: []string{}},
		{name: "Valid all warnings LogTest", args: allWarningArgs, want: allWarningTest, want1: true, want2: []string{}, want3: allExpectedWarns},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := NewLogTest(tt.args.Version, tt.args.RuleID, tt.args.RuleLevel, tt.args.RuleDescription, tt.args.LogFilePath, tt.args.Format, tt.args.Decoder, tt.args.Predecoder, tt.args.TestDescription)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogTest() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NewLogTest() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("NewLogTest() got2 = %v, want %v", got2, tt.want2)
			}
			if !reflect.DeepEqual(got3, tt.want3) {
				t.Errorf("NewLogTest() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func Test_NewLogTestInvalidVersion(t *testing.T) {
	// Completely correct LogTest except for an invalid version
	// which should make NewLogTest return back invalid for the
	// object creation
	logFilePath, err := createTestLog("This is a test log file.")
	if err != nil {
		t.Fatalf("Failed to create test log file")
	}

	version := "9999"
	ruleId := "100"
	ruleLevel := "8"
	ruleDescription := "This is a description"
	// logFilePath already set
	format := "syslog"
	decoder := map[string]string{}
	predecoder := map[string]string{}
	testDescription := "Valid test description."

	expectedTest := new(LogTest)
	expectedTest.Version = version
	expectedTest.RuleID = ruleId
	expectedTest.RuleLevel = ruleLevel
	expectedTest.RuleDescription = ruleDescription
	expectedTest.LogFilePath = logFilePath
	expectedTest.Format = format
	expectedTest.Decoder = decoder
	expectedTest.Predecoder = predecoder
	expectedTest.TestDescription = testDescription

	gotTest, valid, errors, _ := NewLogTest(version, ruleId, ruleLevel, ruleDescription, logFilePath, format, decoder, predecoder, testDescription)

	if !reflect.DeepEqual(gotTest, expectedTest) {
		t.Errorf("NewLogTest() got = %v, want %v", gotTest, expectedTest)
	}

	if valid {
		t.Errorf("Invalid version should render test invalid. Got: %v instead", valid)
	}

	if len(errors) < 1 {
		// Fatal because next test will break
		// if there are no errors
		t.Fatalf("Expected at least one error")
	}

	// Use this approach to search
	// for the error to create flexibility
	// in the future
	hasVersionErr := false
	for _, err := range errors {
		err = strings.ToLower(err)
		if strings.Contains(err, "version") && strings.Contains(err, "invalid") {
			hasVersionErr = true
			break
		}
	}

	if !hasVersionErr {
		t.Errorf("Missing invalid version error")
	}
}

func Test_NewLogTestInvalidRuleID(t *testing.T) {
	// Completely correct LogTest except for an invalid rule id
	// which should make NewLogTest return back invalid for the
	// object creation
	logFilePath, err := createTestLog("This is a test log file.")
	if err != nil {
		t.Fatalf("Failed to create test log file")
	}

	version := "0.1" // First valid version shipped
	ruleId := "9999999999"
	ruleLevel := "8"
	ruleDescription := "This is a description"
	// logFilePath already set
	format := "syslog"
	decoder := map[string]string{}
	predecoder := map[string]string{}
	testDescription := "Valid test description."

	expectedTest := new(LogTest)
	expectedTest.Version = version
	expectedTest.RuleID = ruleId
	expectedTest.RuleLevel = ruleLevel
	expectedTest.RuleDescription = ruleDescription
	expectedTest.LogFilePath = logFilePath
	expectedTest.Format = format
	expectedTest.Decoder = decoder
	expectedTest.Predecoder = predecoder
	expectedTest.TestDescription = testDescription

	gotTest, valid, errors, _ := NewLogTest(version, ruleId, ruleLevel, ruleDescription, logFilePath, format, decoder, predecoder, testDescription)

	if !reflect.DeepEqual(gotTest, expectedTest) {
		t.Errorf("NewLogTest() got = %v, want %v", gotTest, expectedTest)
	}

	if valid {
		t.Errorf("Invalid rule ID should render test invalid. Got: %v instead", valid)
	}

	if len(errors) < 1 {
		// Fatal because next test will break
		// if there are no errors
		t.Fatalf("Expected at least one error")
	}

	// Use this approach to search
	// for the error to create flexibility
	// in the future
	hasVersionErr := false
	for _, err := range errors {
		err = strings.ToLower(err)
		if strings.Contains(err, "rule id") && strings.Contains(err, "invalid") {
			hasVersionErr = true
			break
		}
	}

	if !hasVersionErr {
		t.Errorf("Missing invalid rule id error")
	}
}

func Test_NewLogTestInvalidRuleLevel(t *testing.T) {
	// Completely correct LogTest except for an invalid rule level
	// which should make NewLogTest return back invalid for the
	// object creation
	logFilePath, err := createTestLog("This is a test log file.")
	if err != nil {
		t.Fatalf("Failed to create test log file")
	}

	version := "0.1" // First valid version shipped
	ruleId := "9999"
	ruleLevel := "-15"
	ruleDescription := "This is a description"
	// logFilePath already set
	format := "syslog"
	decoder := map[string]string{}
	predecoder := map[string]string{}
	testDescription := "Valid test description."

	expectedTest := new(LogTest)
	expectedTest.Version = version
	expectedTest.RuleID = ruleId
	expectedTest.RuleLevel = ruleLevel
	expectedTest.RuleDescription = ruleDescription
	expectedTest.LogFilePath = logFilePath
	expectedTest.Format = format
	expectedTest.Decoder = decoder
	expectedTest.Predecoder = predecoder
	expectedTest.TestDescription = testDescription

	gotTest, valid, errors, _ := NewLogTest(version, ruleId, ruleLevel, ruleDescription, logFilePath, format, decoder, predecoder, testDescription)

	if !reflect.DeepEqual(gotTest, expectedTest) {
		t.Errorf("NewLogTest() got = %v, want %v", gotTest, expectedTest)
	}

	if valid {
		t.Errorf("Invalid rule level should render test invalid. Got: %v instead", valid)
	}

	if len(errors) < 1 {
		// Fatal because next test will break
		// if there are no errors
		t.Fatalf("Expected at least one error")
	}

	// Use this approach to search
	// for the error to create flexibility
	// in the future
	hasVersionErr := false
	for _, err := range errors {
		err = strings.ToLower(err)
		if strings.Contains(err, "rule level") && strings.Contains(err, "invalid") {
			hasVersionErr = true
			break
		}
	}

	if !hasVersionErr {
		t.Errorf("Missing invalid rule level error")
	}
}

func Test_NewLogTestInvalidFormat(t *testing.T) {
	// Completely correct LogTest except for an invalid format
	// which should make NewLogTest return back invalid for the
	// object creation
	logFilePath, err := createTestLog("This is a test log file.")
	if err != nil {
		t.Fatalf("Failed to create test log file")
	}

	version := "0.1" // First valid version shipped
	ruleId := "9999"
	ruleLevel := "15"
	ruleDescription := "This is a description"
	// logFilePath already set
	format := "INVALID_FORMAT"
	decoder := map[string]string{}
	predecoder := map[string]string{}
	testDescription := "Valid test description."

	expectedTest := new(LogTest)
	expectedTest.Version = version
	expectedTest.RuleID = ruleId
	expectedTest.RuleLevel = ruleLevel
	expectedTest.RuleDescription = ruleDescription
	expectedTest.LogFilePath = logFilePath
	expectedTest.Format = format
	expectedTest.Decoder = decoder
	expectedTest.Predecoder = predecoder
	expectedTest.TestDescription = testDescription

	gotTest, valid, errors, _ := NewLogTest(version, ruleId, ruleLevel, ruleDescription, logFilePath, format, decoder, predecoder, testDescription)

	if !reflect.DeepEqual(gotTest, expectedTest) {
		t.Errorf("NewLogTest() got = %v, want %v", gotTest, expectedTest)
	}

	if valid {
		t.Errorf("Invalid fomat should render test invalid. Got: %v instead", valid)
	}

	if len(errors) < 1 {
		// Fatal because next test will break
		// if there are no errors
		t.Fatalf("Expected at least one error")
	}

	// Use this approach to search
	// for the error to create flexibility
	// in the future
	hasVersionErr := false
	for _, err := range errors {
		err = strings.ToLower(err)
		if strings.Contains(err, "format") && strings.Contains(err, "invalid") {
			hasVersionErr = true
			break
		}
	}

	if !hasVersionErr {
		t.Errorf("Missing invalid format error")
	}
}
