package main

import (
	"flag"
	"os"
)

func main() {
    flag.String("path", os.Getenv("GODIR"))
}
