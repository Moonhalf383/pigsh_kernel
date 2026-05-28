package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pigsh "github.com/Moonhalf383/pigsh_kernel"
)

// cliIO implements pigsh.IO using stdin/stdout.
type cliIO struct {
	scanner *bufio.Scanner
}

func (c *cliIO) Print(value string) {
	fmt.Println(value)
}

func (c *cliIO) Input(prompt string) string {
	if prompt != "" {
		fmt.Print(prompt + " ")
	}
	c.scanner.Scan()
	return strings.TrimSpace(c.scanner.Text())
}

func main() {
	var source string

	if len(os.Args) > 1 {
		// read from file
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		source = string(data)
	} else {
		// read from stdin
		var buf strings.Builder
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			buf.WriteString(scanner.Text())
			buf.WriteByte('\n')
		}
		source = buf.String()
	}

	io := &cliIO{scanner: bufio.NewScanner(os.Stdin)}
	if err := pigsh.Run(source, io); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
