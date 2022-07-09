package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mdouchement/tac/aio"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	revision = "none"
	date     = "unknown"
)

func main() {
	ctrl := &aio.Controller{
		Logger: &aio.Logger{},
	}

	c := &cobra.Command{
		Use:          "tac",
		Short:        "Tube Amps Calculator",
		Args:         cobra.ExactArgs(1),
		Version:      fmt.Sprintf("%.7s @ %s - %s", revision, date, runtime.Version()),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
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
