package main

import (
	"errors"
	"os"

	"github.com/fatih/color"
	"github.com/s-petr/celeritas"
)

const version = "1.0.0"

var cel celeritas.Celeritas

func main() {
	var message string
	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		exitGracefully(err)
	}

	setup(arg1, arg2)

	switch arg1 {
	case "help":
		showHelp()
	case "new":
		if arg2 == "" {
			exitGracefully(errors.New("please provide a name for the application"))
		}
		doNew(arg2)
	case "version":
		color.Yellow("Application version: %s", version)
	case "make":
		if arg2 == "" {
			exitGracefully(errors.New("make requires a subcommand (migration/model/handler)"))
		}
		if err := doMake(arg2, arg3); err != nil {
			exitGracefully(err)
		}
	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		if err = doMigrate(arg2, arg3); err != nil {
			exitGracefully(err)
		}
		message = "Migrations complete!"
	case "exit":
		exitGracefully(nil)
	default:
		showHelp()
	}
	exitGracefully(nil, message)
}

func validateInput() (string, string, string, error) {
	var arg1, arg2, arg3 string
	argCount := len(os.Args) - 1

	if argCount >= 1 {
		arg1 = os.Args[1]

		if argCount >= 2 {
			arg2 = os.Args[2]

			if argCount >= 3 {
				arg3 = os.Args[3]
			}
		}
	} else {
		color.Red("Error: command required")
		showHelp()
		return "", "", "", errors.New("command required")
	}

	return arg1, arg2, arg3, nil
}

func exitGracefully(err error, msg ...string) {
	message := ""

	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Finished!")
	}

	os.Exit(0)
}
