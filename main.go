package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kinhouse/kh-site/persist"
	"github.com/kinhouse/kh-site/server"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("set the PORT")
	}

	persist, err := persist.NewPersist()
	if err != nil {
		panic(err)
	}

	s := server.BuildServer(persist)
	s.Run(fmt.Sprintf(":%d", port))
}
