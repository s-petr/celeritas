package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
)

var appURL string

func doNew(appName string) {
	appName = strings.ToLower(appName)
	appURL = appName

	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[(len(exploded) - 1)]
	}

	log.Println("App name is", appName)

	color.Green("  Cloning repository...")

	if _, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/s-petr/celeritas-app.git",
		Progress: os.Stdout,
		Depth:    1,
	}); err != nil {
		exitGracefully(err)
	}

	if err := os.RemoveAll(fmt.Sprintf("./%s/.git", appName)); err != nil {
		exitGracefully(err)
	}

	color.Yellow("  Creating .env file...")
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		exitGracefully(err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", cel.RandomString(32))

	if err := copyDataToFile([]byte(env), fmt.Sprintf("%s/.env", appName)); err != nil {
		exitGracefully(err)
	}

	color.Yellow("  Adding Makefile...")

	makeFileExt := "linuxmac"
	if runtime.GOOS == "windows" {
		makeFileExt = "windows"
	}

	if err := copyFileFromTemplate(fmt.Sprintf("templates/Makefile.%s", makeFileExt), fmt.Sprintf("./%s/Makefile", appName)); err != nil {
		exitGracefully(err)
	}

	color.Yellow("  Creating go.mod file...")
	if err := os.Remove("./" + appName + "/go.mod"); err != nil {
		exitGracefully(err)
	}

	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		exitGracefully(err)
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	if err := copyDataToFile([]byte(mod), "./"+appName+"/go.mod"); err != nil {
		exitGracefully(err)
	}

	color.Yellow("  Updating source files...")
	os.Chdir("./" + appName)
	updateSource()

	color.Yellow("  Getting latest version of Celeritas repository...")
	cmd := exec.Command("go", "get", "github.com/s-petr/celeritas")
	if err := cmd.Start(); err != nil {
		exitGracefully(err)
	}

	color.Yellow("  Running go mod tidy...")
	cmd = exec.Command("go", "mod", "tidy")
	if err := cmd.Start(); err != nil {
		exitGracefully(err)
	}

	color.Green("  Done building " + appName)
	color.HiGreen("Go build something awesome")
}
