package main

import (
	"flag"
	"log"

	"github.com/iand/petrel"
)

var (
	addr = flag.String("addr", ":3001", "address to serve on")
	peer = flag.String("peer", "", "peer to bootstrap from")
)

func main() {
	flag.Parse()

	self := petrel.Node{
		ID:   500,
		Addr: *addr,
	}

	ring := &petrel.Ring{
		Owner: self,
	}

	if *peer != "" {

		msg := petrel.Message{
			Purpose: petrel.MsgUser,
			Sender:  self,
			Target: petrel.Node{
				Addr: *peer,
			},
		}

		resp, err := ring.Send(&msg)
		log.Printf("got %v, %v", resp, err)
	}

	log.Printf("listening on %s", self.Addr)
	log.Fatal(petrel.ListenAndServe(*addr))
}
