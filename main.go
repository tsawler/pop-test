package main

import (
	"embed"
	"fmt"
	"github.com/gobuffalo/pop"
	"log"
	"time"
)

const migrationPath = "./migrations"

//go:embed templates
var templateFS embed.FS

func main() {
	// Get a connection to the database. Values are read from config/database.yml
	tx, err := pop.Connect("development")
	if err != nil {
		log.Panic(err)
	}
	defer tx.Close()

	// generate a fizz up/down migration for users table
	up, _ := templateFS.ReadFile("templates/users.postgres.up.fizz")
	down, _ := templateFS.ReadFile("templates/users.postgres.down.fizz")
	err = createMigration(up, down, "users_table", "fizz")
	if err != nil {
		log.Fatal(err)
	}

	// add auto update of updated_at function/trigger up/down migrations as sql
	up, _ = templateFS.ReadFile("templates/users_trigger.postgres.up.sql")
	down, _ = templateFS.ReadFile("templates/users_trigger.postgres.down.sql")
	err = createMigration(up, down, "users_trigger", "sql")
	if err != nil {
		log.Fatal(err)
	}

	// run the migrations
	err = runMigrations(tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

func createMigration(up, down []byte, migrationName, migrationType string) error {
	err := pop.MigrationCreate(migrationPath, migrationName, migrationType, up, down)
	if err != nil {
		return err
	}

	// we have to wait one second so that the migration name is not duplicated
	time.Sleep(1 * time.Second)
	return nil
}

func runMigrations(tx *pop.Connection) error {
	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	// run the migrations
	err = fm.Up()
	if err != nil {
		return err
	}
	return nil
}

func migrateDown(tx *pop.Connection, steps ...int) error {
	step := 1
	if len(steps) > 0 {
		step = steps[0]
	}
	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	// run the migrations
	err = fm.Down(step)
	if err != nil {
		return err
	}
	return nil
}

func migrateReset(tx *pop.Connection) error {
	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	// run the migrations
	err = fm.Reset()
	if err != nil {
		return err
	}
	return nil
}

func steps(tx *pop.Connection, steps int) error {
	return migrateDown(tx, steps)
}
