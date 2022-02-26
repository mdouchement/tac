package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mdouchement/tac/aio"
	"github.com/muesli/coral"
	"gopkg.in/yaml.v2"
)

var (
	revision = "none"
	date     = "unknown"
)

func main() {
	ctrl := &aio.Controller{
		Logger: &aio.Logger{},
	}

	c := &coral.Command{
		Use:          "tac",
		Short:        "Tube Amps Calculator",
		Args:         coral.ExactArgs(1),
		Version:      fmt.Sprintf("build %.7s @ %s - %s", revision, date, runtime.Version()),
		SilenceUsage: true,
		RunE: func(c *coral.Command, args []string) error {
			payload, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			err = yaml.Unmarshal(payload, &ctrl.Config)
			if err != nil {
				return err
			}

			err = ctrl.Execute()
			if err != nil {
				return err
			}

			ctrl.Logger.Verbose("")
			ctrl.Logger.Verbose("====")
			ctrl.Logger.Verbose("")

			ctrl.Logger.Println(ctrl.Results)

			return nil
		},
	}
	c.Flags().BoolVarP(&ctrl.Logger.IsVerbose, "verbose", "V", false, "Verbose logs")

	if err := c.Execute(); err != nil {
		switch {
		case strings.Contains(err.Error(), "unknown shorthand flag"):
			fallthrough
		case strings.Contains(err.Error(), "accepts "):
			c.Println(c.UsageString())
		}
		os.Exit(1)
	}
}
