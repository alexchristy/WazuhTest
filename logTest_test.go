package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
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
		{name: "Invalid negative RuleID -15", args: args{"-15"}, want: false, want1: []string{"Rule ID cannot be less than 0"}, want2: []string{}},
		{name: "Invalid too large RuleID 1234567890", args: args{"1234567890"}, want: false, want1: []string{"Rule ID cannot be greater than 999999"}, want2: []string{}},
		{name: "Invalid non-numeric RuleID hello", args: args{"hello"}, want: false, want1: []string{"Rule ID is not an integer"}, want2: []string{}},
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
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidRuleLevel() got1 = %v, want %v", got1, tt.want1)
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
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("isValidFormat() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("isValidFormat() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Benchmark_isValidFormat_invalidFormats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidFormat(strconv.Itoa(i))
	}
}
