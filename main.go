package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/NarrativeTeam/shawty/handlers"
	"github.com/NarrativeTeam/shawty/storages"
	"github.com/mitchellh/go-homedir"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	dir, _ := homedir.Dir()
	storage := &storages.Filesystem{}
	err := storage.Init(filepath.Join(dir, "shawty"))
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", handlers.MainHandler(storage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
