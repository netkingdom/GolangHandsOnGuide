package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
)

type config struct {
	numTimes int
	path     string
}

var errPosArgsSpecified = errors.New("Positional arguments specified")
var usageString = fmt.Sprintf(`Usage: %s <integer> [-h|--help]

A greeter application which prints the name you entered <integer> number 
of times.
`, os.Args[0])

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

	fs.StringVar(&c.path, "o", fmt.Sprintf("%s", "./index.html"), "test")

	err := fs.Parse(args) // flag 옵션 검사
	if err != nil {
		return c, err
	}

	if fs.NArg() != 0 { //플래그 옵션이 파싱된 이후에 주어진 위치 인수의 개수 return
		return c, errPosArgsSpecified
	}
	return c, nil
}

func validateArgs(c config) error {
	if !(c.numTimes > 0) {
		return errors.New("Must specify a number greater than 0")
	}
	return nil
}

func greetUser(c config, name string, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, msg)
	}
}

func greetHTML(c config, name string) error {
	if c.path == "" {
		return nil
	}
	var contents string
	const tpl = `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
				<title>TEST</title>
			</head>
			<body>
				<div>{{ . }}</div>
			</body>
		</html>`

	f, err := os.Create(c.path)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	t, err := template.New("Greet").Parse(tpl)
	if err != nil {
		return err
	}

	for i := 0; i < c.numTimes; i++ {
		contents += fmt.Sprintf("<h1>Nice to meet you %s</h1>\n", name)
	}

	err = t.Execute(w, template.HTML(contents))
	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

func runCmd(r io.Reader, w io.Writer, c config) error {
	name, err := getName(r, w)
	if err != nil {
		return err
	}

	greetUser(c, name, w)
	err = greetHTML(c, name)
	if err != nil {
		return err
	}

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
