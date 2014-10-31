package web

import (
	"fmt"
	"github.com/LapisBlue/Tar/cli"
	"github.com/ogier/pflag"
	"os"
	"path/filepath"
)

const (
	serverConfig = "lapitar.json"
)

func Run(name string, args []string) int {
	flags := pflag.NewFlagSet(name, pflag.ContinueOnError)

	dir := flags.StringP("dir", "d", ".", "The folder to save all files in.")
	config := flags.StringP("config", "c", serverConfig, "The configuration file used to configure the server.")

	cli.FlagUsage(name, flags)

	if len(args) >= 1 && args[0] == "help" {
		flags.Usage()
		return 1
	}

	if flags.Parse(args) != nil {
		return 1
	}

	if *dir != "." && filepath.Dir(*config) == "." {
		*config = filepath.Join(*dir, *config)
	}

	// Load the configuration
	fmt.Println("Loading configuration from:", *config)
	conf, exit := loadConfig(*config)
	if conf == nil {
		return exit
	}

	return 0
}

func loadConfig(path string) (conf *config, exit int) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		// Create a new configuration file
		file, err = os.Create(path)
		if err != nil {
			exit = cli.PrintError(err, "Failed to create configuration file")
			return
		}

		defer file.Close()

		err = writeConfig(file, defaultConfig())
		if err == nil {
			fmt.Println("Created configuration:", path)
		} else {
			exit = cli.PrintError(err, "Failed to write configuration")
		}
	} else if err != nil {
		exit = cli.PrintError(err, "Failed to open configuration file")
	} else {
		defer file.Close()

		// Read the configuration from the file
		conf, err = parseConfig(file)
		if err != nil {
			exit = cli.PrintError(err, "Failed to parse configuration")
		}
	}

	return
}