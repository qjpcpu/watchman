package alfred

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
    "log"
    "os"
    "router"
    "strings"
    "syscall"
    "time"
)

const TimeFormat = "2006-01-02 15:04:05"

var pool *WatcherPool

func init() {
    pool = initPool()
    pool.boot()
    // Start router
    router.Start(router.DefaultBuilder())
    cli, err := router.NewRouterCli(router.SYS_ID, router.DefaultBuilder().SocketFunc)
    if err != nil {
        log.Fatalln(err)
    }
    // Biding emitter
    distributer := &Distributer{cli}
    pool.emitter = distributer
    // Start receive from router
    go func() {
        for {
            if msg, err := distributer.PullRequest(); err == nil {
                pool.Signal <- msg
            }
        }
    }()
}

type Distributer struct {
    *router.RouterCli
}

// The request message must like:
//{
//    "Event":0,
//    "FileName":"+/path/to/file",  or "FileName":"-/path/to/file"  // +/- means add or remove watch
//    ....
//}
func (em *Distributer) PullRequest() (map[string]string, error) {
    str, err := em.Read()
    if err != nil {
        log.Println("Pull request:", err)
        return nil, err
    }
    m, err := router.ParseMessage(str)
    if err != nil {
        log.Println("Pull request:", err)
        return nil, err
    }
    if m.Event != 0 {
        return nil, errors.New("Invalid request")
    }
    msg := make(map[string]string)
    if strings.HasPrefix(m.FileName, "+") {
        msg["ACTION"] = "ADD"
        msg["PATH"] = strings.TrimLeft(m.FileName, "+")
        return msg, nil
    } else if strings.HasPrefix(m.FileName, "-") {
        msg["ACTION"] = "REMOVE"
        msg["PATH"] = strings.TrimLeft(m.FileName, "-")
        return msg, nil
    }
    return nil, errors.New("Shouldn't come here")
}

// The event operation is:
//{
//    "Mask":0,
//    "Name":"FAIL:/path/to/file" or "Name":"SUCCESS:/path/to/file"
//}
func (em *Distributer) Eject(env *inotify.Event, t time.Time) {
    //log.Println("alfred.go", env)
    if env.Mask == 0x0 {
        m := router.Message{
            Event:    0x0,
            FileName: env.Name,
        }
        em.Write(m.String())
    } else {
        m1 := router.Message{
            Event:    env.Mask,
            FileName: env.Name,
        }
        buildMsg(env.Name, &m1)
        em.Write(m1.String())
    }
}

func buildMsg(path string, msg *router.Message) {
    if fi, err := os.Stat(path); err == nil {
        msg.Size = fi.Size()
        if t, ok := fi.Sys().(*syscall.Stat_t); ok {
            msg.Inode = t.Ino
            msg.AccessTime = time.Unix(t.Atim.Unix()).Format(TimeFormat)
            msg.ChangeTime = time.Unix(t.Ctim.Unix()).Format(TimeFormat)
            msg.ModifyTime = time.Unix(t.Mtim.Unix()).Format(TimeFormat)
        } else {
            log.Println("Can't get %v details by syscall", path)
        }
    }
}
