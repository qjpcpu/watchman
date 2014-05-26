package watchman

import (
    "errors"
    "router"
    "strconv"
    "strings"
    "time"
)

type Watchman struct {
    paths  map[string]uint32
    client *router.RouterCli
}

// Initialize a new watchman
func NewWatchman() (*Watchman, error) {
    c, err := router.NewRouterCli(strconv.Itoa(time.Now().Nanosecond()), router.DefaultBuilder().SocketFunc)
    if err != nil {
        return nil, err
    }
    return &Watchman{make(map[string]uint32), c}, nil
}

// Add a file path to watch list, specify the events as you need
func (man *Watchman) WatchPath(path string, events uint32) error {
    if _, ok := man.paths[path]; ok {
        man.paths[path] = events & IN_ALL_EVENTS
        return nil
    }
    if len(path) > 1 && strings.HasSuffix(path, "/") {
        path = strings.TrimRight(path, "/")
    }
    m := router.Message{
        Event:    0x0,
        FileName: "+" + path,
    }
    err := man.client.Write(m.String())
    if err != nil {
        return err
    }
    man.paths[path] = events & IN_ALL_EVENTS
    return nil
}

// Stop watching a path
func (man *Watchman) ForgetPath(path string) error {
    if _, ok := man.paths[path]; !ok {
        return nil
    }
    if len(path) > 1 && strings.HasSuffix(path, "/") {
        path = strings.TrimRight(path, "/")
    }
    m := router.Message{
        Event:    0x0,
        FileName: "-" + path,
    }
    err := man.client.Write(m.String())
    if err != nil {
        return err
    }
    delete(man.paths, path)
    return nil
}

// Fetch an event of watching list, if there's no event available the function would blocked
func (man *Watchman) PullEvent() (router.Message, error) {
    raw, err := man.client.Read()
    if err != nil {
        return router.Message{}, err
    }
    m, err := router.ParseMessage(raw)
    if err != nil {
        return router.Message{}, err
    }
    fn := ""
    for name, _ := range man.paths {
        if strings.HasPrefix(m.FileName, fn) {
            fn = name
        }
    }
    if m.Event == 0x0 || fn == "" || man.paths[fn]&m.Event == 0 {
        return router.Message{}, errors.New("You dont' need it.")
    }
    m.Event = m.Event & IN_ALL_EVENTS & man.paths[fn]
    return m, nil
}

// Get all watching files
func (man *Watchman) CheckPathList() []string {
    list := make([]string, len(man.paths))
    i := 0
    for k, _ := range man.paths {
        list[i] = k
        i += 1
    }
    return list
}

// Stop watching and release resources
func (man *Watchman) Release() {
    man.client.Close()
}
