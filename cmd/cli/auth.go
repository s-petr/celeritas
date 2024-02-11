package main

import (
	"github.com/fatih/color"
)

func doAuth() error {
	dbType := cel.DB.DataType

	tx, err := cel.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	// fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	// upFile := cel.RootPath + "/migrations/" + fileName + ".up.sql"
	// downFile := cel.RootPath + "/migrations/" + fileName + ".down.sql"

	// if err := copyFileFromTemplate("templates/migrations/auth_tables."+
	// 	dbType+".sql", upFile); err != nil {
	// 	exitGracefully(err)
	// }

	upBytes, err := templateFS.ReadFile("templates/migrations/auth_tables." + dbType + ".sql")
	if err != nil {
		exitGracefully(err)
	}

	stmt := `
	DROP TABLE IF EXISTS users CASCADE;
	DROP TABLE IF EXISTS tokens CASCADE;
	DROP TABLE IF EXISTS remember_tokens CASCADE;
	`

	downBytes := []byte(stmt)

	if err := cel.CreatePopMigration(upBytes, downBytes, "auth", "sql"); err != nil {
		exitGracefully(err)
	}

	// if err := copyDataToFile([]byte(stmt), downFile); err != nil {
	// 	exitGracefully(err)
	// }

	// if err := doMigrate("up", ""); err != nil {
	// 	exitGracefully(err)
	// }

	if err := cel.RunPopMigrations(tx); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/data/user.go.txt",
		cel.RootPath+"/data/user.go"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/data/token.go.txt",
		cel.RootPath+"/data/token.go"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/data/remember-token.go.txt",
		cel.RootPath+"/data/remember-token.go"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/middleware/auth.go.txt",
		cel.RootPath+"/middleware/auth.go"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/middleware/auth-token.go.txt",
		cel.RootPath+"/middleware/auth-token.go"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/middleware/remember.go.txt",
		cel.RootPath+"/middleware/remember.go"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/handlers/auth-handlers.go.txt",
		cel.RootPath+"/handlers/auth-handlers.go"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/mail/password-reset.html.tmpl",
		cel.RootPath+"/mail/password-reset.html.tmpl"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/mail/password-reset.text.tmpl",
		cel.RootPath+"/mail/password-reset.text.tmpl"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/views/login.jet",
		cel.RootPath+"/views/login.jet"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/views/forgot.jet",
		cel.RootPath+"/views/forgot.jet"); err != nil {
		exitGracefully(err)
	}

	if err := copyFileFromTemplate("templates/views/reset-password.jet",
		cel.RootPath+"/views/reset-password.jet"); err != nil {
		exitGracefully(err)
	}

	color.Yellow("  - users, tokens, and remember-tokens migrations created and executed")
	color.Yellow("  - users and token models created")
	color.Yellow("  - auth middleware created")
	color.Yellow("")
	color.Yellow("Please add user and token models in data/models.go")
	color.Yellow("Add the appropriate middleware to your routes")
	return nil
}
