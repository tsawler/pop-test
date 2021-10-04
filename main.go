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
	tx, err := pop.Connect("development")
	if err != nil {
		log.Panic(err)
	}
	defer tx.Close()

	// generate a fizz up/down migration for users table
	usersUp, _ := templateFS.ReadFile("templates/users.postgres.up.fizz")
	usersDown, _ := templateFS.ReadFile("templates/users.postgres.down.fizz")
	err = createMigration(usersUp, usersDown, "users_table", "fizz")
	if err != nil {
		log.Fatal(err)
	}

	// add auto update of updated_at function/trigger, and down migration
	triggerUp, _ := templateFS.ReadFile("templates/users_trigger.postgres.up.sql")
	triggerDown, _ := templateFS.ReadFile("templates/users_trigger.postgres.down.sql")
	err = createMigration(triggerUp, triggerDown, "users_trigger", "sql")
	if err != nil {
		log.Fatal(err)
	}

	// test table
	whateverUp, _ := templateFS.ReadFile("templates/whatever.postgres.up.fizz")
	whateverDown, _ := templateFS.ReadFile("templates/whatever.postgres.down.fizz")
	err = createMigration(whateverUp, whateverDown, "whatever", "fizz")
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

	// strangely, we have to wait one second so that the migration name is not duplicated
	time.Sleep(1 * time.Second)
	return nil
}

func runMigrations(tx *pop.Connection) error {
	fm, err := pop.NewFileMigrator("./migrations", tx)
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
