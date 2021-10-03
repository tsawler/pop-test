package main

import (
	"embed"
	"github.com/gobuffalo/pop"
	"log"
)


//go:embed templates
var templateFS embed.FS

func main() {
	tx, err := pop.Connect("development")
	if err != nil {
		log.Panic(err)
	}
	defer tx.Close()

	usersUp, _ := templateFS.ReadFile("templates/users.postgres.up.fizz")
	usersDown, _ := templateFS.ReadFile("templates/users.postgres.down.fizz")

	err = pop.MigrationCreate("./migrations", "users_table", "fizz", usersUp, usersDown)
	if err != nil {
		log.Fatal(err)
	}

	triggerUp, _ := templateFS.ReadFile("templates/users_trigger.postgres.up.sql")
	triggerDown, _ := templateFS.ReadFile("templates/users_trigger.postgres.down.sql")

	err = pop.MigrationCreate("./migrations", "users_trigger", "sql", triggerUp, triggerDown)
	if err != nil {
		log.Fatal(err)
	}
	fm, err := pop.NewFileMigrator("./migrations", tx)
	if err != nil {
		log.Fatal(err)
	}

	err = fm.Up()
	if err != nil {
		log.Fatal(err)
	}
}