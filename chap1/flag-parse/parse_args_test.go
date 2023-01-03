package main

import (
	"bytes"
	"errors"
	"testing"
)

type testConfig struct {
	args     []string
	err      error
	numTimes int
}

// TODO - 이전에 정의한 testConfig 구조체 삽입
func TestParseArgs(t *testing.T) {
	tests := []testConfig{
		{
			args:     []string{"-h"},
			err:      nil,
			numTimes: 0,
		},
		{
			args:     []string{"-n", "10"},
			err:      nil,
			numTimes: 10,
		},
		{
			args:     []string{"-n", "abc"},
			err:      errors.New("invalid value \"abc\" for flag -n: parse error"),
			numTimes: 0,
		},
		{
			args:     []string{"-n", "1", "foo"},
			err:      errors.New("Positional arguments specified"),
			numTimes: 1,
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		c, err := parseArgs(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n", tc.err, err)
		}
		if c.numTimes != tc.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.numTimes, c.numTimes)
		}
		byteBuf.Reset()
	}
}
