package petrel

import (
	"net"
)

type Client struct {
	conn net.Conn
}

func (c *Client) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// func (c *Client) Maint(ttl int, events []Event) error {
// 	msg := &Message{
// 		Type:   MessageMaint,
// 		TTL:    ttl,
// 		Events: events,
// 	}

// 	proto := new(ProtocolV1)
// 	return proto.WriteMessage(c.conn, msg)
// }

// func (c *Client) Bootstrap() (*BootstrapResponse, error) {
// 	msg := &Message{
// 		Type: MessageBootstrap,
// 		TTL:  0,
// 	}

// 	proto := new(ProtocolV1)
// 	err := proto.WriteMessage(c.conn, msg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	buf := bufio.NewReader(io.LimitReader(c.conn, 4096))

// err = proto.ReadResponseLine(buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return proto.ReadBootstrap(buf)
// }

// func (c *Client) Ping() error {
// 	msg := &Message{
// 		Type: MessagePing,
// 		TTL:  0,
// 	}

// 	proto := new(ProtocolV1)
// 	err := proto.WriteMessage(c.conn, msg)
// 	if err != nil {
// 		return err
// 	}

// 	buf := bufio.NewReader(io.LimitReader(c.conn, 4096))
// 	return proto.ReadResponseLine(buf)
// }
