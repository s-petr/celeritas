package celeritas

import (
	"log"

	"github.com/gobuffalo/pop"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (c *Celeritas) PopConnect() (*pop.Connection, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}
	return tx, nil

}

func (c *Celeritas) getFileMigrator(tx *pop.Connection) (*pop.FileMigrator, error) {
	var migrationPath = c.RootPath + "/migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return nil, err
	}
	return &fm, nil
}

func (c *Celeritas) CreatePopMigration(up, down []byte, migrationName, migrationType string) error {
	var migrationPath = c.RootPath + "/migrations"
	return pop.MigrationCreate(migrationPath, migrationName, migrationType, up, down)
}

func (c *Celeritas) RunPopMigrations(tx *pop.Connection) error {
	fm, err := c.getFileMigrator(tx)
	if err != nil {
		return err
	}

	return fm.Up()
}

func (c *Celeritas) PopMigrateDown(tx *pop.Connection, steps ...int) error {
	step := 1
	if len(steps) > 0 {
		step = steps[0]
	}

	fm, err := c.getFileMigrator(tx)
	if err != nil {
		return err
	}

	return fm.Down(step)
}

func (c *Celeritas) PopMigrateReset(tx *pop.Connection) error {
	fm, err := c.getFileMigrator(tx)
	if err != nil {
		return err
	}

	return fm.Reset()
}

func (c *Celeritas) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("error running migration:", err)
		return err
	}
	return nil
}

func (c *Celeritas) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Down()
}

func (c *Celeritas) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Steps(n)
}

func (c *Celeritas) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Force(-1)
}
