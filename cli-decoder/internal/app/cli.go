package app

import (
	"bufio"
	"cli-decoder/internal/decode"
	"fmt"
	"io"
)

type CLI struct {
	decoder decode.Decoder
	in      *bufio.Scanner
	out     io.Writer
}

func NewCLI(decoder decode.Decoder, in io.Reader, out io.Writer) *CLI {
	return &CLI{
		decoder: decoder,
		in:      bufio.NewScanner(in),
		out:     out,
	}
}

func (c CLI) Run() {
	fmt.Fprintln(c.out, "Application run...")
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprint(c.out, "Recovered. Error:\n", r)
		}
	}()

	for {
		if err := c.decoder.Save(c.readLine()); err != nil {
			fmt.Fprint(c.out, "Error: ", err)
		}
	}
}

func (c CLI) readLine() string {
	c.in.Scan()
	return c.in.Text()
}
