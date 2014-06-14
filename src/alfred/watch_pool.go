package alfred

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
    . "mlog"
    "path/filepath"
    "reflect"
    "strings"
    "time"
)

type Emitter interface {
    Eject(*inotify.Event, time.Time)
}

// WatcherPool control all the watchers
type WatcherPool struct {
    Table   map[string]*alfredWatcher // The Table shows the paths and its according watcher
    List    []*alfredWatcher          // The List includes all the alfredwatchers
    Signal  chan map[string]string    // The Singal is a channel, used by communication
    emitter Emitter
    counter map[string]int
}

// Initialize a watch pool, this is a private package function
func initPool() *WatcherPool {
    return &WatcherPool{
        make(map[string]*alfredWatcher),
        []*alfredWatcher{},
        make(chan map[string]string),
        nil,
        make(map[string]int),
    }
}

// Get the listen file list
func (wp *WatcherPool) FileList() []string {
    list := make([]string, len(wp.Table))
    i := 0
    for k, _ := range wp.Table {
        list[i] = k
        i += 1
    }
    return list
}

// Pick up a watcher whose listen list is not full, then  add the path to its list
func (wp *WatcherPool) Attach(path string) error {
    if _, ok := wp.Table[path]; ok {
        wp.counter[path] += 1
        return nil
    }
    var w *alfredWatcher
    for _, v := range wp.List {
        if v.Size() < MAX_PATH_PER_WATCHER {
            w = v
            break
        }
    }
    if w == nil {
        if len(wp.List) >= MAX_WATCHER {
            return errors.New("To much watchers.")
        }
        w = newAlfredWatcher()
        wp.List = append(wp.List, w)
        Log.Debug("Create new watcher for ", path)
    }
    if err := w.AddWatch(path); err != nil {
        return err
    } else {
        wp.Table[path] = w
        wp.counter[path] += 1
    }
    Log.Debug(path + " is under watching...")
    return nil
}
func (wp *WatcherPool) Dettach(path string) error {
    if w, ok := wp.Table[path]; !ok {
        return nil
    } else {
        wp.counter[path] -= 1
        if wp.counter[path] == 0 {
            err := w.RemoveWatch(path)
            if err != nil {
                Log.Debug(err.Error())
                return err
            }
            delete(wp.Table, path)
            delete(wp.counter, path)
            Log.Debugf("Remove %v from watching list.", path)
        } else {
            Log.Debugf("Remove a reference to %v from watching list.", path)
        }
    }
    return nil
}

func (wp *WatcherPool) GetDefaultPaths() []string {
    return []string{}
}
func (wp *WatcherPool) schedule() {
    var cases []reflect.SelectCase
    flush := true
    for {
        if flush {
            cases = make([]reflect.SelectCase, len(wp.List)+1)
            cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(wp.Signal)}
            for i, ch := range wp.List {
                cases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.watcher.Event)}
            }
            flush = false
        }
        chosen, value, _ := reflect.Select(cases)
        if chosen == 0 {
            msg := value.Interface().(map[string]string)
            wp.handleMessage(msg)
            flush = true
        } else {
            ev := value.Interface().(*inotify.Event)
            if wp.emitter != nil {
                go wp.emitter.Eject(ev, time.Now())
            }
        }
    }
}
func (wp *WatcherPool) handleMessage(msg map[string]string) {
    var err error
    var action string
    if path := msg["PATH"]; msg["ACTION"] == "ADD" {
        action = "+"
        err = wp.Attach(path)
    } else if msg["ACTION"] == "REMOVE" {
        action = "-"
        err = wp.Dettach(path)
    }
    if wp.emitter == nil {
        return
    }
    if err != nil {
        env := &inotify.Event{0, 0, "FAIL:" + action + msg["PATH"]}
        go wp.emitter.Eject(env, time.Now())
    } else {
        env := &inotify.Event{0, 0, "SUCCESS:" + action + msg["PATH"]}
        go wp.emitter.Eject(env, time.Now())
    }

}
func (wp *WatcherPool) boot() {
    for _, fn := range wp.GetDefaultPaths() {
        if err := wp.Attach(fn); err != nil {
            Log.Debug(err.Error())
        }
    }
    go wp.schedule()
}
func (wp *WatcherPool) shutdown() {
    for fn, _ := range wp.Table {
        wp.Dettach(fn)
    }
    for _, w := range wp.List {
        w.Release()
    }
}
func (wp *WatcherPool) triggerPaths(event_path string) (pkeys []string) {
    event_fn := event_path
    if strings.HasPrefix(event_path, "SUCCESS:") {
        event_fn = event_path[9:]
    } else if strings.HasPrefix(event_path, "FAIL:") {
        event_fn = event_path[6:]
    }
    if strings.HasSuffix(event_fn, "/") {
        event_fn = strings.TrimRight(event_fn, "/")
    }
    if _, ok := wp.Table[event_fn]; ok {
        pkeys = append(pkeys, event_fn)
    }
    dir, _ := filepath.Split(event_fn)
    if strings.HasSuffix(dir, "/") {
        dir = strings.TrimRight(dir, "/")
    }
    if _, ok := wp.Table[dir]; ok {
        pkeys = append(pkeys, dir)
    }
    return
}
