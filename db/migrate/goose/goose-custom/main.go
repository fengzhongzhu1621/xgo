package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	assets "github.com/fengzhongzhu1621/xgo/db/migrate/goose"
	_ "github.com/fengzhongzhu1621/xgo/db/migrate/goose/migrations"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"trpc.group/trpc-go/tnet/log"
)

// Usage of goose:
//
//		-dir string
//	  		directory with migration files (default ".")
var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	flags.Usage = usage

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("goose: failed to parse flags: %v", err)
	}
	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}
	if len(args) < 3 {
		flags.Usage()
		return
	}

	dbstring, command := args[1], args[2]
	log.Infof("args[1] = %s, dbString = %s, command = %s", args[0], dbstring, command)

	db, err := goose.OpenDBWithDriver(args[0], dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Errorf("goose: failed to close DB: %v", err)
		}
	}()

	goose.SetVerbose(true)

	policy := backoff.NewExponentialBackOff()
	policy.MaxElapsedTime = 30 * time.Second
	err = backoff.Retry(func() error {
		return db.PingContext(context.Background())
	}, policy)
	if err != nil {
		log.Fatalf("failed to initialize database connection: %w", err)
	}

	goose.SetBaseFS(assets.EmbedMigrations)

	currentVersion, err := goose.GetDBVersion(db)
	if err != nil {
		log.Fatalf("failed to get db version: %w", err)
	}
	log.Infof("current db version = %d", currentVersion)

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}

	log.Info("migration done")
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: goose-custom COMMAND
Examples:
    goose-custom status
`
	usageCommands = `
Commands:
   up                   Migrate the DB to the most recent version available
   up-by-one            Migrate the DB up by 1
   up-to VERSION        Migrate the DB to a specific VERSION
   down                 Roll back the version by 1
   down-to VERSION      Roll back to a specific VERSION
   redo                 Re-run the latest migration
   reset                Roll back all migrations
   status               Dump the migration status for the current DB
   version              Print the current version of the database
   create NAME [sql|go] Creates new migration file with the current timestamp
   fix                  Apply sequential ordering to migrations`
)
