package alfred

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
    "log"
    "reflect"
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
}

// Initialize a watch pool, this is a private package function
func initPool() *WatcherPool {
    return &WatcherPool{
        make(map[string]*alfredWatcher),
        []*alfredWatcher{},
        make(chan map[string]string),
        nil,
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
    if _, ok := wp.Table["path"]; ok {
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
        log.Println("Create new watcher for ", path)
    }
    if err := w.AddWatch(path); err != nil {
        return err
    } else {
        wp.Table[path] = w
    }
    log.Println(path + " is under watching...")
    return nil
}
func (wp *WatcherPool) Dettach(path string) error {
    if w, ok := wp.Table["path"]; !ok {
        log.Println(path + " has not been watched.")
        return nil
    } else {
        err := w.RemoveWatch(path)
        if err != nil {
            return err
        }
        delete(wp.Table, path)
    }
    return nil
}

func (wp *WatcherPool) GetDefaultPaths() []string {
    return []string{"/tmp", "/home/work", "/home/work/tmp/"}
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
    if path := msg["PATH"]; msg["ACTION"] == "ADD" {
        err = wp.Attach(path)
    } else if msg["ACTION"] == "REMOVE" {
        err = wp.Dettach(path)
    }
    if wp.emitter == nil {
        return
    }
    if err != nil {
        env := &inotify.Event{0, 0, "FAIL:" + msg["PATH"]}
        go wp.emitter.Eject(env, time.Now())
    } else {
        env := &inotify.Event{0, 0, "SUCCESS:" + msg["PATH"]}
        go wp.emitter.Eject(env, time.Now())
    }

}
func (wp *WatcherPool) boot() {
    for _, fn := range wp.GetDefaultPaths() {
        if err := wp.Attach(fn); err != nil {
            log.Println(err)
        }
    }
    go wp.schedule()
}
