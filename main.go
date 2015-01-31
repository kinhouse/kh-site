package main

import (
	"github.com/kinhouse/kh-site/persist"
	"github.com/kinhouse/kh-site/server"
)

func main() {
	persist, err := persist.NewPersist()
	if err != nil {
		panic(err)
	}
	server := server.BuildServer(persist)
	server.Run()
}
