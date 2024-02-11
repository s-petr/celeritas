package main

import (
	"errors"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

func doMake(arg2, arg3, arg4 string) error {
	switch arg2 {
	case "key":
		rnd := cel.RandomString(32)
		color.Yellow("32 character encryption key: %s", rnd)
	case "migration":
		checkForDB()

		// dbType := cel.DB.DataType

		if arg3 == "" {
			exitGracefully(errors.New("you need to provide a name for the migration"))
		}

		migrationType := "fizz"
		var up, down string

		if arg4 == "fizz" || arg4 == "" {
			upBytes, _ := templateFS.ReadFile("templates/migrations/migration_up.fizz")
			downBytes, _ := templateFS.ReadFile("templates/migrations/migration_down.fizz")

			up = string(upBytes)
			down = string(downBytes)
		} else {
			migrationType = "sql"
		}

		if err := cel.CreatePopMigration([]byte(up), []byte(down),
			arg3, migrationType); err != nil {
			exitGracefully(err)
		}

		// fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		// upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		// downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		// if err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql",
		// 	upFile); err != nil {
		// 	exitGracefully(err)
		// }

		// if err := copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql",
		// 	downFile); err != nil {
		// 	exitGracefully(err)
		// }
	case "auth":
		if err := doAuth(); err != nil {
			exitGracefully(err)
		}
	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("please provide a name for the handler"))
		}

		fileName := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToGoCamel(arg3))

		if err := os.WriteFile(fileName, []byte(handler), 0644); err != nil {
			exitGracefully(err)
		}
	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("please provide a name for the model"))
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		model := string(data)

		plur := pluralize.NewClient()

		var modelName = arg3
		var tableName = arg3

		if plur.IsPlural(arg3) {
			modelName = plur.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(plur.Plural(arg3))
		}

		fileName := cel.RootPath + "/data/" + strings.ToLower(modelName) + ".go"

		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToGoCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

		if err := copyDataToFile([]byte(model), fileName); err != nil {
			exitGracefully(err)
		}
	case "mail":
		if arg3 == "" {
			exitGracefully(errors.New("please provide a name for the email template"))
		}
		htmlMail := cel.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainTextMail := cel.RootPath + "/mail/" + strings.ToLower(arg3) + ".text.tmpl"

		if err := copyFileFromTemplate("templates/mail/mail.html.tmpl", htmlMail); err != nil {
			exitGracefully(err)
		}

		if err := copyFileFromTemplate("templates/mail/mail.text.tmpl", plainTextMail); err != nil {
			exitGracefully(err)
		}

	case "session":
		if err := doSessionTable(); err != nil {
			exitGracefully(err)
		}
	default:
		showHelp()
	}

	return nil
}
