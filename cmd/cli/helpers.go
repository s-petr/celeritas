package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup(arg1, arg2 string) {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		if err := godotenv.Load(); err != nil {
			exitGracefully(err)
		}
	}

	path, err := os.Getwd()
	if err != nil {
		exitGracefully(err)
	}

	cel.RootPath = path
	cel.DB.DataType = os.Getenv("DATABASE_TYPE")
}

func getDSN() string {
	dbType := cel.DB.DataType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		return dsn
	}
	return "mysql://" + cel.BuildDSN()
}

func checkForDB() {
	dbType := cel.DB.DataType

	if dbType == "" {
		exitGracefully(errors.New("no database connection provided in .env"))
	}

	if !fileExists(cel.RootPath + "/config/database.yml") {
		exitGracefully(errors.New("missing config file (/config/database.yml)"))
	}
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}

	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			exitGracefully(err)
		}

		newContents := strings.Replace(string(read), "myapp", appURL, -1)

		if err := os.WriteFile(path, []byte(newContents), 0); err != nil {
			exitGracefully(err)
		}
	}
	return nil
}

func updateSource() {
	if err := filepath.Walk(".", updateSourceFiles); err != nil {
		exitGracefully(err)
	}
}

func showHelp() {
	color.Yellow(`Available commands:

help                            - show the help commands
down                            - set the server into maintenance mode
up                              - take the server out of maintenance mode
version                         - print application version
new                             - create new application from built-in template
migrate                         - runs all up migrations that have not been run previously
migrate down                    - reverses the most recent migration
migrate reset                   - runs all down migrations in reverse order, and then all up migrations
make migration <name> <format>  - creates two new up and down migrations in the migrations folder; format = fizz (default) or sql
make auth                       - creates and runs migrations for users table
make handler <name>             - creates a stub handler in the handlers directory
make model <name>               - creates a new model in the data directory
make session                    - creates a table in the database as a session store
make mail                       - creates two starter mail templates in hte mail directory
	`)
}
