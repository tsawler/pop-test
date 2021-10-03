package main

import (
	"embed"
	"fmt"
	"github.com/gobuffalo/pop"
	"log"
)

const migrations = "./migrations"

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

	err = pop.MigrationCreate(migrations, "users_table", "fizz", usersUp, usersDown)
	if err != nil {
		log.Fatal(err)
	}

	// add auto update of updated_at function/trigger, and down migration
	triggerUp, _ := templateFS.ReadFile("templates/users_trigger.postgres.up.sql")
	triggerDown, _ := templateFS.ReadFile("templates/users_trigger.postgres.down.sql")

	err = pop.MigrationCreate(migrations, "users_trigger", "sql", triggerUp, triggerDown)
	if err != nil {
		log.Fatal(err)
	}

	// create a file migrator
	fm, err := pop.NewFileMigrator("./migrations", tx)
	if err != nil {
		log.Fatal(err)
	}

	// run the migrations
	err = fm.Up()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}
