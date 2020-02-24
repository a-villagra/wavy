package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"math"

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
	buf, err := d.FullPCMBuffer()
	if err != nil {
		errLog.Printf("Failed to get PCM buffer: %v", err)
		return err
	}

	format := buf.Format
	outLog.Printf("File info: Bit rate: %d, Format: %v", buf.SourceBitDepth, format)

	nSamples := 0
	for i := 0; i < len(buf.Data); i++ {
		buf.Data[i] = int(float64(buf.Data[i]) * math.Cos(float64(nSamples) * 0.1))
		nSamples++
	}
	outLog.Printf("Write %d samples", nSamples)

	of, err := os.Create("out.wav")
	if err != nil {
		errLog.Printf("Failed to create output file: %v", err)
		return err
	}
	defer of.Close()
	e := wav.NewEncoder(
		of,
		format.SampleRate,
		int(buf.SourceBitDepth),
		format.NumChannels,
		int(d.WavAudioFormat),
	)
	err = e.Write(buf)
	if err != nil {
		return err
	}
	e.Close()

	return nil
}
