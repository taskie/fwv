package cli

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
	"github.com/taskie/fwv"
	"github.com/taskie/osplus"
)

var (
	log      = logrus.New()
	version  = fwv.Version
	revision = ""
)

type Options struct {
	Mode           string
	ReverseMode    bool   `short:"r" long:"reverse" description:"reverse mode"`
	NoWidth        bool   `short:"W" long:"noWidth" description:"NOT use char width"`
	EaaHalfwidth   bool   `short:"E" long:"eaaHalfWidth" env:"FWV_EAA_HALF_WIDTH" description:"treat East Asian Ambiguous as half width"`
	Color          func() `long:"color" description:"colorize output"`
	NoColor        func() `long:"noColor" description:"NOT colorize output"`
	OutputFilePath string `short:"o" long:"output" description:"output file path"`
	Whitespaces    string `short:"s" long:"whitespaces" description:"characters treated as whitespace"`
	Delimiter      string `short:"d" long:"delimiter" description:"delimiter used for FWV output"`
	Verbose        bool   `short:"v" long:"verbose" description:"show verbose output"`
	Version        bool   `short:"V" long:"version" description:"show version"`
}

func Main() {
	var opts Options

	outFd := os.Stdout.Fd()
	colored := isatty.IsTerminal(outFd) || isatty.IsCygwinTerminal(outFd)
	opts.Color = func() {
		colored = true
	}
	opts.NoColor = func() {
		colored = false
	}

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
	if opts.Mode != "" {
		mode = opts.Mode
	}
	eastAsianAmbiguousWidth := 2
	if opts.EaaHalfwidth {
		eastAsianAmbiguousWidth = 1
	}

	opener := osplus.NewOpener()
	w, commit, err := opener.CreateTempFileWithDestination(opts.OutputFilePath, "", "fwv-")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	r, err := opener.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	conv := fwv.NewConverter(w, r, mode)
	conv.UseWidth = !opts.NoWidth
	conv.EastAsianAmbiguousWidth = eastAsianAmbiguousWidth
	if opts.Whitespaces != "" {
		conv.Whitespaces = opts.Whitespaces
	}
	conv.Delimiter = opts.Delimiter
	if opts.OutputFilePath == "" {
		conv.Colored = colored
	}
	err = conv.Convert()
	if err != nil {
		log.Fatal(err)
	}
	commit(true)
	return
}
