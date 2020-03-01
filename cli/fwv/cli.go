package fwv

import (
	"io"
	"os"
	"path/filepath"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/taskie/fwv"
	"github.com/taskie/ose"
	"github.com/taskie/ose/coli"
	"go.uber.org/zap"
)

var configFile string
var config Config

const CommandName = "fwv"

var Command *cobra.Command

func init() {
	Command = NewCommand(coli.NewColiInThisWorld())
}

func Main() {
	Command.Execute()
}

func NewCommand(cl *coli.Coli) *cobra.Command {
	cmd := &cobra.Command{
		Use:  CommandName,
		Args: cobra.RangeArgs(0, 2),
		Run:  cl.WrapRun(run),
	}
	cl.Prepare(cmd)

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", `config file (default "`+CommandName+`.yml")`)
	cmd.Flags().StringP("from-type", "f", "", "convert from [fwv|csv]")
	cmd.Flags().StringP("to-type", "t", "", "convert to [fwv|csv]")
	cmd.Flags().BoolP("no-width", "W", false, "NOT use char width")
	cmd.Flags().BoolP("eaa-half-width", "E", false, "treat East Asian Ambiguous as half width")
	cmd.Flags().BoolP("show-column-ranges", "r", false, "show column ranges")
	cmd.Flags().BoolP("no-trim", "T", false, "NOT trim whitespaces")
	cmd.Flags().BoolP("color", "C", false, "colorize output")
	cmd.Flags().BoolP("no-color", "M", false, "NOT colorize output (monochrome)")
	cmd.Flags().StringP("whitespaces", "s", " ", "characters treated as whitespace")
	cmd.Flags().StringP("delimiter", "d", " ", "delimiter used for FWV output")

	cl.BindFlags(cmd.Flags(), []string{
		"from-type", "to-type", "no-width", "eaa-half-width", "show-column-ranges", "no-trim",
		"color", "no-color", "whitespaces", "delimiter",
	})
	cl.Viper().BindEnv("no_color")
	return cmd
}

type Config struct {
	Input, Output, FromType, ToType, Whitespaces, Delimiter         string
	NoWidth, EaaHalfWidth, ShowColumnRanges, NoTrim, Color, NoColor bool
}

func run(cl *coli.Coli, cmd *cobra.Command, args []string) {
	v := cl.Viper()
	log := zap.L()
	if v.GetBool("version") {
		cmd.Println(fwv.Version)
		return
	}
	err := v.Unmarshal(&config)
	if err != nil {
		log.Fatal("can't unmarshal config", zap.Error(err))
	}

	input := ""
	output := ""
	switch len(args) {
	case 0:
		break
	case 1:
		input = args[0]
	case 2:
		input = args[0]
		output = args[1]
	default:
		log.Fatal("invalid arguments", zap.Strings("arguments", args[2:]))
	}

	fromType := config.FromType
	if fromType == "" && filepath.Ext(input) == ".csv" {
		fromType = "csv"
	}
	toType := config.ToType
	if toType == "" && filepath.Ext(output) == ".csv" {
		toType = "csv"
	}

	eastAsianAmbiguousWidth := 2
	if config.EaaHalfWidth {
		eastAsianAmbiguousWidth = 1
	}

	colored := false
	if output == "" || output == "-" {
		colored = isatty.IsTerminal(os.Stdout.Fd())
	}
	if config.NoColor && !config.Color {
		colored = false
	} else if !config.NoColor && config.Color {
		colored = true
	}

	opener := ose.NewOpenerInThisWorld()
	r, err := opener.Open(input)
	if err != nil {
		log.Fatal("can't open", zap.Error(err))
	}
	defer r.Close()
	_, err = opener.CreateTempFile("", CommandName, output, func(f io.WriteCloser) (bool, error) {
		conv := fwv.NewConverter(f, r, fromType, toType)
		conv.UseWidth = !config.NoWidth
		conv.EastAsianAmbiguousWidth = eastAsianAmbiguousWidth
		conv.Whitespaces = config.Whitespaces
		conv.Delimiter = config.Delimiter
		conv.Colored = colored
		conv.ShowColumnRanges = config.ShowColumnRanges
		conv.NoTrim = config.NoTrim

		err = conv.Convert()
		if err != nil {
			log.Fatal("can't convert", zap.Error(err))
		}
		return true, nil
	})
	if err != nil {
		log.Fatal("can't create file", zap.Error(err))
	}
}
