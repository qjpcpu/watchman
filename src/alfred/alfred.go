package alfred

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
    "log"
    "router"
    "strings"
)

var pool *WatcherPool

func init() {
    pool = initPool()
    pool.boot()
    // Start router
    router.Start(router.DefaultPolicy)
    cli, err := router.NewRouterCli(router.SYS_ID)
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
    return nil, errors.New("Shoudn't come here")
}
func (em *Distributer) Eject(env *inotify.Event, info string) {
    if env != nil {
        log.Println(env)
    } else {
        log.Println(info)
    }
}
