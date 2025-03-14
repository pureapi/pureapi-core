package main

import (
	"flag"
	"fmt"
)

// Run with: go run -tags=example . --ex=database
func main() {
	example := flag.String("ex", "", "Which example to run (server|middlware|database)")
	flag.Parse()

	switch *example {
	case "server":
		Server()
	case "middleware":
		RunMiddleware()
	case "database":
		RunDatabase()
	default:
		fmt.Println("Usage: go run . --ex=emitter|middlware|database")
	}
}
