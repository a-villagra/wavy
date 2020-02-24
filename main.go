package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"

	"github.com/go-audio/wav"
)

var (
	FlagInputFile string
)

func main() {
	err := Run(os.Stdout, os.Stderr)
	if err != nil {
		os.Stderr.WriteString("Failed to run\n")
		os.Exit(-1)
	}
}

func Run(stdout, stderr io.Writer) error {
	flag.StringVar(&FlagInputFile, "input", "", "the input file (must be a wav file)")
	flag.Parse()

	errLog := log.New(stderr, "ERR ", log.Ltime)
	outLog := log.New(stdout, "INFO ", log.Ltime)

	if FlagInputFile == "" {
		errLog.Print("Input file is mandatory")
		return errors.New("Bad args")
	}

	outLog.Printf("Opening file %s", FlagInputFile)
	f, err := os.Open(FlagInputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	d := wav.NewDecoder(f)
	duration, err := d.Duration()
	if err != nil {
		return err
	}

	outLog.Printf("Original duration: %v", duration)

	return nil
}
