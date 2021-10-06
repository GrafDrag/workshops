package app

import (
	"bufio"
	"cli-decoder/internal/decode"
	"fmt"
	"io"
)

type App struct {
	decoder decode.Decoder
	in      *bufio.Scanner
	out     io.Writer
}

func NewCLI(decoder decode.Decoder, in io.Reader, out io.Writer) *App {
	return &App{
		decoder: decoder,
		in:      bufio.NewScanner(in),
		out:     out,
	}
}

func (c App) Run() {
	fmt.Fprintln(c.out, "Application run...")
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprint(c.out, "Recovered. Error:\n", r)
		}
	}()

	for {
		if err := c.decoder.SetHash(c.readLine()); err != nil {
			fmt.Fprint(c.out, "Error: ", err)
		}
	}
}

func (c App) readLine() string {
	c.in.Scan()
	return c.in.Text()
}
