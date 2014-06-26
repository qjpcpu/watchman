package alfred

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
    "fmt"
    . "mlog"
    "os"
    "path/filepath"
    "router"
    "strings"
    "syscall"
    "time"
    //   "utils"
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
    distributer := &Distributer{cli, make(map[string]time.Time), make(map[int]uint64)}
    pool.emitter = distributer
    // Start receive from router
    go func() {
        for {
            if msg, err := distributer.PullRequest(); err == nil {
                pool.Signal <- msg
            }
        }
    }()
    Log.Info("Alfred: Startup.")
}
func Shutdown() {
    pool.shutdown()
    Log.Info("Alfred: Shutdown.")
}

type Distributer struct {
    *router.RouterCli
    memo     map[string]time.Time
    freqctrl map[int]uint64
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

func (em *Distributer) ctrlDelay(t time.Time) int {
    grade := 5
    delay := 1
    if em.freqctrl[t.Second()] == 0 {
        em.freqctrl[(t.Second()-grade)%60] = 0
    }
    em.freqctrl[t.Second()] += 1
    var all uint64
    for i := 0; i < grade; i++ {
        all += em.freqctrl[(t.Second()-i)%60]
    }
    switch {
    case all > 5000:
        delay = 30
    case all > 3000:
        delay = 10
    default:
        delay = 1
    }
    //use freqctrl[60] as debug tag, the if block(4 lines below) can be deleted.
    if em.freqctrl[60] != uint64(delay) {
        em.freqctrl[60] = uint64(delay)
        Log.Infof("Alfred: got %v notify in last 5 seconds, adjust event eject cycle to %v seconds", all, delay)
    }

    return delay
}

func (em *Distributer) passby(env *inotify.Event, t time.Time) (can_eject bool) {
    delay := em.ctrlDelay(t)
    key := fmt.Sprintf("%s:%v", filepath.Dir(env.Name), env.Mask)
    if last, ok := em.memo[key]; !ok {
        em.memo[key] = t
        can_eject = true
    } else {
        if t.Before(last.Add(time.Duration(delay) * time.Second)) {
            can_eject = false
        } else {
            can_eject = true
            em.memo[key] = t
        }
    }
    if env.Mask == 0x0 {
        can_eject = true
    }
    return
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
        if em.passby(env, t) {
            em.Write(to, m.String())
        }
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
