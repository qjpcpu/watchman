package router

import (
    "net"
)

type RouterCli struct {
    conn net.Conn
    id   string
}

func NewRouterCli(uid string) (*RouterCli, error) {
    cli := &RouterCli{
        id: uid,
    }
    c, err := net.Dial("unix", socketpath())
    if err != nil {
        return cli, err
    }
    // Handshake
    if err := writeString(c, uid); err != nil {
        return cli, err
    }
    if res, err := readString(c); err != nil || res != "connected" {
        return cli, err
    }
    cli.conn = c
    return cli, nil
}

func (rc *RouterCli) Write(msg string) error {
    return writeString(rc.conn, msg)
}
func (rc *RouterCli) Read() (string, error) {
    return readString(rc.conn)
}

func (rc *RouterCli) Close() error {
    return rc.conn.Close()
}
