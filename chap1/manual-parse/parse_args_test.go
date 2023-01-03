package main

import (
	"errors"
	"testing"
)

type testConfig struct {
	args []string
	err  error
	config
}

// TODO - 이전에 정의한 testConfig 구조체 삽입
func TestParseArgs(t *testing.T) {
	tests := []testConfig{
		{
			args: []string{"-h"},
			err:  nil,
			config: config{
				numTimes:   0,
				printUsage: true,
			},
		},
		{
			args: []string{"10"},
			err:  nil,
			config: config{
				numTimes:   10,
				printUsage: false,
			},
		},
		{
			args: []string{"abc"},
			err:  errors.New("strconv.Atoi: parsing \"abc\": invalid syntax"),
			config: config{
				numTimes:   0,
				printUsage: false,
			},
		},
		{
			args: []string{"1", "foo"},
			err:  errors.New("Invalid number of arguments"),
			config: config{
				numTimes:   0,
				printUsage: false,
			},
		},
	}

	for _, tc := range tests {
		c, err := parseArgs(tc.args)
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n", tc.err, err)
		}
		if tc.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if c.printUsage != tc.printUsage {
			t.Errorf("Expected printUsage to be: %v, got: %v\n", tc.config.printUsage, c.printUsage)
		}
		if c.numTimes != tc.config.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.config.numTimes, c.numTimes)
		}
	}
}
