package petrel

import (
	"fmt"
	"io"
)

const (
	MsgSys = iota
	MsgUser
)

type Messager interface {
	Send(*Message) (*Message, error)
}

type Message struct {
	ID      uint64
	Purpose int
	Sender  Node
	Target  Node
	TTL     int
	Body    io.ReadCloser
}

type Codec interface {
	Decode(io.Reader, *Message) error
	Encode(io.Writer, *Message) error
}

type TextCodec struct{}

func (t *TextCodec) Encode(w io.Writer, msg *Message) error {
	_, err := fmt.Fprintf(w, "%d %d %d %s %d\n", msg.ID, msg.Purpose, msg.Sender.ID, msg.Sender.Addr, msg.TTL)
	return err
}

func (t *TextCodec) Decode(r io.Reader, msg *Message) error {
	_, err := fmt.Fscanf(r, "%d %d %d %s %d\n", &msg.ID, &msg.Purpose, &msg.Sender.ID, &msg.Sender.Addr, &msg.TTL)
	if err != nil {
		return err
	}

	return err
}
