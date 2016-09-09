package main

import (
	"flag"
	"log"

	"github.com/iand/petrel"
)

var (
	peer = flag.String("-peer", "", "peer to bootstrap from")
)

func main() {
	log.Fatal(petrel.ListenAndServe(":3001"))
}
