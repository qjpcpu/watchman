package router

import (
    "bytes"
    "encoding/binary"
    "log"
    "net"
    "os"
    "strings"
)

// Router server can route every message to other connections by policy
type RouterServ struct {
    conns        map[string]net.Conn
    policy       ForwardPolicy
    center_queue chan map[string]string
    lock         chan int
    listener     net.Listener
}

// Forward policy
type ForwardPolicy func(id, body, toid string) bool

// Get the unix socket path
type SocketPath func() string

// Contains forwrd policy and socket path getter
type Builder struct {
    ForwardFunc ForwardPolicy
    SocketFunc  SocketPath
}

// Accept new connection and add to server's client connections
func (ss *RouterServ) serv() {
    go ss.forwardMsg()
    for {
        fd, err := ss.listener.Accept()
        if err != nil {
            log.Println(err)
        } else {
            go ss.addClient(fd)
        }
    }
}

func writeString(c net.Conn, msg string) error {
    buf := new(bytes.Buffer)
    data := []byte(msg)
    err := binary.Write(buf, binary.LittleEndian, int16(len(data)))
    if err != nil {
        return err
    }
    buf.Write(data)
    if _, err := c.Write(buf.Bytes()); err != nil {
        return err
    }
    return nil
}
func readString(c net.Conn) (string, error) {
    size_data := make([]byte, 2)
    if n, err := c.Read(size_data); err != nil || n != 2 {
        return "", err
    }
    var size int16
    if err := binary.Read(bytes.NewReader(size_data), binary.LittleEndian, &size); err != nil {
        return "", err
    }
    data := make([]byte, size)
    if _, err := c.Read(data); err != nil {
        return "", err
    } else {
        return string(data), nil
    }

}

// If the first message from client is id, and the id doesn't exist, the connection is ok.
func (ss *RouterServ) addClient(c net.Conn) {
    id, err := readString(c)
    if err != nil {
        log.Println("Handshake error", err)
        c.Close()
        return
    }
    if _, ok := ss.conns[id]; ok {
        log.Println("id " + id + " already exists!")
        c.Close()
        return
    } else {
        ss.lock <- 1
        writeString(c, "connected")
        ss.conns[id] = c
        <-ss.lock
    }
    for {
        msg := make(map[string]string)
        body, err := readString(c)
        if err == nil {
            msg["id"] = id
            msg["body"] = body
            ss.center_queue <- msg
        } else if err.Error() == "EOF" {
            ss.removeClient(id)
            log.Printf("Remove %v from router\n", id)
            break
        }
    }
}
func (ss *RouterServ) removeClient(id ...string) {
    ss.lock <- 1
    for _, key := range id {
        if c, ok := ss.conns[key]; ok {
            c.Close()
            delete(ss.conns, key)
        }
    }
    <-ss.lock
}

// forward messages to other connections by policy
func (ss *RouterServ) forwardMsg() {
    for {
        msg := <-ss.center_queue
        if ss.policy != nil {
            var broken []string
            // lock
            ss.lock <- 1
            for k, v := range ss.conns {
                if ok := ss.policy(msg["id"], msg["body"], k); ok {
                    err := writeString(v, msg["body"])
                    if err != nil {
                        log.Printf("%v -> %v [%v] Error:%v\n", msg["id"], k, msg["body"], err)
                        if strings.HasSuffix(err.Error(), "broken pipe") {
                            broken = append(broken, k)
                        }
                    }
                }
            }
            //unlock
            <-ss.lock
            // If has broken pipes, remove them
            if len(broken) > 0 {
                ss.removeClient(broken...)
            }
        } else {
            log.Println("Empty forward policy, drop message!")
        }
    }
}

func Start(builder Builder) (*RouterServ, error) {
    path := builder.SocketFunc()
    serv := &RouterServ{make(map[string]net.Conn), builder.ForwardFunc, make(chan map[string]string, 100), make(chan int, 1), nil}
    os.Remove(path)
    l, err := net.Listen("unix", path)
    if err != nil {
        log.Println("start message server err", err)
        return serv, err
    } else {
        log.Println("Message server started!")
        serv.listener = l
    }
    go serv.serv()
    return serv, nil
}
