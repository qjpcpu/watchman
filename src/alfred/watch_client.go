package alfred

import (
    "container/list"
    "errors"
    "path/filepath"
    "strings"
    "time"
)

type Watchman struct {
    paths    map[string]uint32
    messages *list.List
}

// Initialize a new watchman
func NewWatchman() (*Watchman, error) {
    return &Watchman{paths: make(map[string]uint32), messages: list.New()}, nil
}

// Add a file path to watch list, specify the events as you need
func (man *Watchman) WatchPath(path string, events uint32) error {
    if _, ok := man.paths[path]; ok {
        man.paths[path] = events
        return nil
    }
    if len(path) > 1 && strings.HasSuffix(path, "/") {
        path = strings.TrimRight(path, "/")
    }
    GetManager().Attach(man, path)
    man.paths[path] = events
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
    GetManager().Dettach(man, path)
    delete(man.paths, path)
    return nil
}

func (man *Watchman) MessageRecieved(m Message) {
    fn := ""
    for name, _ := range man.paths {
        if inWatch(m.FileName, name) {
            fn = name
            break
        }
    }
    if _, ok := man.paths[fn]; ok && m.Event == 0x0 {
        return
    }
    if m.Event == 0x0 || fn == "" || man.paths[fn]&m.Event == 0 {
        return
    }
    man.messages.PushFront(m)
}

// Fetch an event of watching list, if there's no event available the function would blocked
func (man *Watchman) PullEvent() (Message, error) {
    for {
        if ele := man.messages.Back(); ele != nil {
            value := man.messages.Remove(ele)
            return value.(Message), nil
        } else {
            time.Sleep(time.Millisecond * 100)
        }
    }
    return Message{}, errors.New("shouldn't come here")
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
}
func inWatch(event_path, fn string) bool {
    event_fn := event_path
    if strings.HasPrefix(event_path, "SUCCESS:") {
        event_fn = event_path[9:]
    } else if strings.HasPrefix(event_path, "FAIL:") {
        event_fn = event_path[6:]
    }
    if strings.HasSuffix(event_fn, "/") {
        event_fn = strings.TrimRight(event_fn, "/")
    }
    if strings.HasSuffix(fn, "/") {
        fn = strings.TrimRight(fn, "/")
    }
    if event_fn == fn {
        return true
    }
    dir, _ := filepath.Split(event_fn)
    if strings.HasSuffix(dir, "/") {
        dir = strings.TrimRight(dir, "/")
    }
    if fn == dir {
        return true
    }
    return false
}
