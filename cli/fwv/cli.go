package fwv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-isatty"

	"github.com/taskie/fwv"

	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taskie/jc"
	"github.com/taskie/osplus"
)

type Config struct {
	Input, Output, FromType, ToType, Whitespaces, Delimiter, LogLevel string
	NoWidth, EaaHalfWidth, ShowColumnRanges, Color, NoColor           bool
}

var configFile string
var config Config
var (
	verbose, debug, version bool
)

const CommandName = "fwv"

func init() {
	Command.PersistentFlags().StringVarP(&configFile, "config", "c", "", `config file (default "`+CommandName+`.yml")`)
	Command.Flags().StringP("fromType", "f", "", "convert from [fwv|csv]")
	Command.Flags().StringP("toType", "t", "", "convert to [fwv|csv]")
	Command.Flags().BoolP("noWidth", "W", false, "NOT use char width")
	Command.Flags().BoolP("eaaHalfWidth", "E", false, "treat East Asian Ambiguous as half width")
	Command.Flags().BoolP("showColumnRanges", "r", false, "show column ranges")
	Command.Flags().BoolP("color", "C", false, "colorize output")
	Command.Flags().BoolP("noColor", "M", false, "NOT colorize output (monochrome)")
	Command.Flags().StringP("whitespaces", "s", " ", "characters treated as whitespace")
	Command.Flags().StringP("delimiter", "d", "", "delimiter used for FWV output")
	Command.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	Command.Flags().BoolVarP(&debug, "debug", "g", false, "debug output")
	Command.Flags().BoolVarP(&version, "version", "V", false, "show Version")

	for _, s := range []string{"fromType", "toType", "noWidth", "eaaHalfWidth", "showColumnRanges", "color", "noColor", "whitespaces", "delimiter"} {
		viper.BindPFlag(s, Command.Flags().Lookup(s))
	}

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(CommandName)
		conf, err := osplus.GetXdgConfigHome()
		if err != nil {
			log.Info(err)
		} else {
			viper.AddConfigPath(filepath.Join(conf, CommandName))
		}
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix(CommandName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Warn(err)
	}
}

func Main() {
	Command.Execute()
}

var Command = &cobra.Command{
	Use:  CommandName,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		err := run(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func run(cmd *cobra.Command, args []string) error {
	if version {
		fmt.Println(jc.Version)
		return nil
	}
	if config.LogLevel != "" {
		lv, err := log.ParseLevel(config.LogLevel)
		if err != nil {
			log.Warn(err)
		} else {
			log.SetLevel(lv)
		}
	}
	if debug {
		if viper.ConfigFileUsed() != "" {
			log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		}
		log.Debug(pp.Sprint(config))
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
		return fmt.Errorf("invalid arguments: %v", args[2:])
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

	opener := osplus.NewOpener()
	r, err := opener.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, commit, err := opener.CreateTempFileWithDestination(output, "", CommandName+"-")
	if err != nil {
		return err
	}
	defer w.Close()

	colored := false
	if output == "" || output == "-" {
		colored = isatty.IsTerminal(os.Stdout.Fd())
	}
	if config.NoColor && !config.Color {
		colored = false
	} else if !config.NoColor && config.Color {
		colored = true
	}

	conv := fwv.NewConverter(w, r, fromType, toType)
	conv.UseWidth = !config.NoWidth
	conv.EastAsianAmbiguousWidth = eastAsianAmbiguousWidth
	conv.Whitespaces = config.Whitespaces
	conv.Delimiter = config.Delimiter
	conv.Colored = colored
	conv.ShowColumnRanges = config.ShowColumnRanges

	err = conv.Convert()
	if err != nil {
		return err
	}
	commit(true)
	return nil
}
