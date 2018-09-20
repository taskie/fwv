package cli

import (
	"bufio"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/taskie/fwv"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	version  = fwv.Version
	revision = ""
)

type Options struct {
	ReverseMode  bool `short:"r" long:"reverse" description:"reverse mode"`
	NoWidth      bool `short:"W" long:"noWidth" description:"NOT use char width"`
	EaaHalfwidth bool `short:"E" long:"eaaHalfWidth" env:"FWV_EAA_HALF_WIDTH" description:"treat East Asian Ambiguous as half width"`
	// NoColor        bool   `long:"no-color" env:"NO_COLOR" description:"NOT colorize output"`
	OutputFilePath string `short:"o" long:"output" description:"output file path"`
	Verbose        bool   `short:"v" long:"verbose" description:"show verbose output"`
	Version        bool   `short:"V" long:"version" description:"show version"`
}

func move(dst string, src string) error {
	err := os.Rename(src, dst)
	if err != nil {
		ifs, err := os.Open(src)
		if err != nil {
			return err
		}
		ofs, err := os.Create(dst)
		if err != nil {
			return err
		}
		_, err = io.Copy(ofs, ifs)
		if err != nil {
			return err
		}
		ofs.Close()
		ifs.Close()
		err = os.Remove(src)
		if err != nil {
			return err
		}
	}
	return nil
}

func Main() {
	var opts Options
	args, err := flags.ParseArgs(&opts, os.Args)
	if opts.Version {
		if opts.Verbose {
			fmt.Println("Version:    ", version)
			if revision != "" {
				fmt.Println("Revision:   ", revision)
			}
		} else {
			fmt.Println(version)
		}
		os.Exit(0)
	}
	if len(args) >= 3 {
		log.Fatal("you can't specify multiple args.")
	}
	if err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	mode := "c2f"
	if opts.ReverseMode {
		mode = "f2c"
	}
	eastAsianAmbiguousWidth := 2
	if opts.EaaHalfwidth {
		eastAsianAmbiguousWidth = 1
	}

	app := fwv.NewApplication(mode)
	app.UseWidth = !opts.NoWidth
	app.EastAsianAmbiguousWidth = eastAsianAmbiguousWidth

	var r *bufio.Reader
	if len(args) == 1 || args[1] == "-" {
		r = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		r = bufio.NewReader(file)
	}

	var w *bufio.Writer
	if opts.OutputFilePath == "" || opts.OutputFilePath == "-" {
		w = bufio.NewWriter(os.Stdout)
	} else {
		file, err := ioutil.TempFile("", "fwv")
		tmpName := file.Name()
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			file.Close()
			err := move(opts.OutputFilePath, tmpName)
			if err != nil {
				log.Fatal(err)
			}
		}()
		w = bufio.NewWriter(file)
	}
	defer w.Flush()

	err = app.Run(r, w)
	if err != nil {
		log.Fatal(err)
	}
}
