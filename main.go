package main

import (
	"os"
	"strconv"

	"github.com/kinhouse/kh-site/persist"
	"github.com/kinhouse/kh-site/server"
)

func main() {
	persist, err := persist.NewPersist()
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("set the PORT")
	}

	server := server.BuildServer(persist)
	server.Run(port)
}
