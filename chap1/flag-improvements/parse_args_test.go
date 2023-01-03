package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		args []string
		config
		output string
		err    error
	}{
		{
			// -h 전달했을 때
			args: []string{"-h"},
			output: `A greeter application which prints the name you entered a specified number of times.

		Usage of greeter: <options> [name]
Options: 
  -n int
    	Number of times to greet
`,
			err:    errors.New("flag: help requested"),
			config: config{numTimes: 0},
		},
		{
			// n을 전달할때
			// 정상
			args:   []string{"-n", "10"},
			err:    nil,
			config: config{numTimes: 10},
		},
		{
			// 비정상 정상
			args:   []string{"-n", "abc"},
			err:    errors.New("invalid value \"abc\" for flag -n: parse error"),
			config: config{numTimes: 0},
		},
		{
			// 위치인수로 이름을 지정할 때
			// 정상
			args:   []string{"-n", "1", "KIM"},
			err:    nil,
			config: config{numTimes: 1, name: "KIM"},
		},
		{
			// 위치인수로 이름을 지정할 때
			// 정상
			args:   []string{"-n", "1", "KIM", "JUN"},
			err:    errors.New("More than one positional arguments specified"),
			config: config{numTimes: 1},
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		c, err := parseArgs(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n", tc.err, err)
		}
		if c.numTimes != tc.numTimes {
			t.Errorf("Exptedted numTimes to be: %v, got: %v",
				tc.numTimes, c.numTimes)
		}

		gotMsg := byteBuf.String()
		if len(tc.output) != 0 && gotMsg != tc.output {
			t.Errorf("Expected stdout message to be: %#v, Got: %#v\n", tc.output, gotMsg)
		}

		byteBuf.Reset()
	}

}
