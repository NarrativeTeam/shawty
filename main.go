package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/NarrativeTeam/shawty/handlers"
	"github.com/NarrativeTeam/shawty/storages"
	"github.com/getsentry/raven-go"
	"github.com/mitchellh/go-homedir"
)

var migrate bool

func init() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	flag.BoolVar(&migrate, "migrate", false, "Specify to migrate database")
}

func main() {
	flag.Parse()

	// Seed the randomizer for the token-generation
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	dir, _ := homedir.Dir()
	postgresHost := os.Getenv("POSTGRES_HOST")
	var storage storages.IStorage
	if postgresHost != "" {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbName := os.Getenv("POSTGRES_DB")
		ssl := os.Getenv("POSTGRES_SSL") != "false"
		pg, err := storages.NewPostgres(postgresHost, dbName, user, password, ssl)
		if err != nil {
			//Could not correct to postgres
			raven.CaptureError(err, nil, nil)
			panic(err)
		}
		storage = pg
		if migrate {
			pg.Setup()
			return
		}
	} else {
		str := &storages.Filesystem{}
		err := str.Init(filepath.Join(dir, "shawty"))
		if err != nil {
			raven.CaptureError(err, nil, nil)
		}
		storage = str
	}

	http.HandleFunc("/", raven.RecoveryHandler(handlers.MainHandler(storage)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		raven.CaptureError(err, nil, nil)
	}
}
