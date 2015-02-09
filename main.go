package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kinhouse/kh-site/server"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("set the PORT")
	}

	s := server.BuildServer()
	s.Run(fmt.Sprintf(":%d", port))
}
