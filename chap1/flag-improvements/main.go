package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	numTimes int
	name     string
}

var errPosArgsSpecified = errors.New("More than one positional arguments specified")

func getName(r io.Reader, w io.Writer) (string, error) {
	msg := "Your name please? Press the Enter key when done.\n"
	fmt.Fprintf(w, msg)
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You didn't enter your name")
	}
	return name, nil
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}

	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)     // flag 객체 생성
	fs.SetOutput(w)                                            // FlagSet의 진단 메세지 혹은 출력메세지를 작성하는데 사용할 Writer Set
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet") // int 타입 플래그 옵션

	fs.Usage = func() {
		var usageString = `A greeter application which prints the name you entered a specified number of times.

		Usage of %s: <options> [name]`

		fmt.Fprintf(w, usageString, fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	err := fs.Parse(args) // flag 옵션 검사
	if err != nil {
		return c, err
	}

	if fs.NArg() > 1 { //플래그 옵션이 파싱된 이후에 주어진 위치 인수의 개수 return
		return c, errPosArgsSpecified
	}

	if fs.NArg() == 1 { //플래그 옵션이 파싱된 이후에 주어진 위치 인수의 개수 return
		c.name = fs.Arg(0)
	}
	return c, nil
}

func validateArgs(c config) error {
	if !(c.numTimes > 0) {
		return errors.New("Must specify a number greater than 0")
	}
	return nil
}

func greetUser(c config, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", c.name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, msg)
	}
}
func runCmd(rd io.Reader, w io.Writer, c config) error {
	var err error
	if len(c.name) == 0 {
		c.name, err = getName(rd, w)
		if err != nil {
			return err
		}

	}

	greetUser(c, w)

	return nil
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		if errors.Is(err, errPosArgsSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		os.Exit(1)
	}

	err = validateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

}
