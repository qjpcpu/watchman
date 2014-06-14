package alfred

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
    . "mlog"
    "os"
    "router"
    "strings"
    "syscall"
    "time"
)

const TimeFormat = "2006-01-02 15:04:05"

var pool *WatcherPool

func Boot() {
    pool = initPool()
    pool.boot()
    // Start router
    cli := router.NewRouterCli(router.SYS_ID, router.DefaultBuildClient)
    cli.Subscribe(router.SYS_ID)
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
    Log.Debug("Alfred Startup.")
}
func Shutdown() {
    pool.shutdown()
    Log.Debug("Alfred Shutdown.")
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
        Log.Debug("Pull request:", err)
        return nil, err
    }
    m, err := router.ParseMessage(str)
    if err != nil {
        Log.Debug("Pull request:", err)
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
    var m router.Message
    if env.Mask == 0x0 {
        m = router.Message{
            Event:    0x0,
            FileName: env.Name,
        }
    } else {
        m = router.Message{
            Event:    env.Mask,
            FileName: env.Name,
        }
        buildMsg(env.Name, &m)
    }
    to_list := pool.triggerPaths(env.Name)
    for _, to := range to_list {
        em.Write(to, m.String())
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
            Log.Debug("Can't get %v details by syscall", path)
        }
    }
}
