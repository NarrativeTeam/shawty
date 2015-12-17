package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/NarrativeTeam/shawty/handlers"
	"github.com/NarrativeTeam/shawty/storages"
	"github.com/mitchellh/go-homedir"
)

func main() {
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
		pg, err := storages.NewPostgres(postgresHost, user, password, dbName, ssl)
		if err != nil {
			//Could not correct to postgres
			log.Panic(err)
		}
		storage = pg
	} else {
		str := &storages.Filesystem{}
		err := str.Init(filepath.Join(dir, "shawty"))
		if err != nil {
			log.Fatal(err)
		}
		storage = str
	}

	http.Handle("/", handlers.MainHandler(storage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
