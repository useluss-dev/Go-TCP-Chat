package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type Conn struct {
	conn   net.Conn
	reader *bufio.Reader
}

func Dial(addr string) *Conn {
	c, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	return &Conn{conn: c, reader: bufio.NewReader(c)}
}

func (c *Conn) Send(line string) error {
	_, err := fmt.Fprintf(c.conn, "%s\n", line)
	return err
}

func (c *Conn) ReadLine() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line, nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
