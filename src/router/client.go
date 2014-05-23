package router

import (
    "log"
    "net"
    "os"
)

type RouterCli struct {
    conn        net.Conn
    id string
}

func NewRouterCli(uid string) *RouterCli,error{
    cli:=&RouterCli{
        id:uid,
    }
    cli.conn,err=net.Dial("unix",socketpath())
    if err!=nil{
        return cli,err
    }
    // Handshake
    _,err:=cli.conn.Write("")
}
// Accept new connection and add to server's client connections
func (ss *RouterServ) serv() {
    go ss.forwordMsg()
    for {
        fd, err := ss.listener.Accept()
        if err != nil {
            log.Println(err)
        } else {
            go ss.addClient(fd)
        }
    }
}

// If the first message from client is {"id":"/path/to","body":"anything"}, and the id doesn't exist, the connection is ok.
func (ss *RouterServ) addClient(c net.Conn) {
    data := make([]byte, 512)
    n, err := c.Read(data)
    if err != nil {
        log.Println("Handshake error", err)
        c.Close()
        return
    }
    id := string(data[0:n])
    if _, ok := ss.conns[id]; ok {
        log.Println("id " + id + " already exists!")
        c.Close()
        return
    } else {
        ss.conns[id] = c
    }
    msg := make(map[string]string)
    for {
        data = make([]byte, 512)
        if n, err = c.Read(data); err == nil {
            msg["id"] = id
            msg["body"] = string(data)
            ss.center_queue <- msg
        }
    }
}
func (ss *RouterServ) removeClient(c net.Conn) {
}

// forword messages to other connections by policy
func (ss *RouterServ) forwordMsg() {
    for {
        msg := <-ss.center_queue
        if ss.policy != nil {
            list := ss.policy(msg["id"], msg["body"])
            tbl := ss.conns
            for _, key := range list {
                if c, ok := tbl[key]; ok {
                    _, err := c.Write([]byte(msg["body"]))
                    if err != nil {
                        log.Println("Forword ", msg["id"], msg["body"], " error. ", err)
                    }
                }
            }
        } else {
            log.Println("Empty forword policy, drop message!")
        }
    }
}

func socketpath() string {
    return "/tmp/router.socket"
}
func Start(forword func(id, body string) []string) (*RouterServ, error) {
    path := socketpath()
    serv := &RouterServ{make(map[string]net.Conn), forword, make(chan map[string]string, 100), nil}
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
