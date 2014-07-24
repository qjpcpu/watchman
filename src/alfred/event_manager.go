package alfred

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
    "github.com/qjpcpu/logger"
    "time"
)

type Emitter interface {
    Eject(*inotify.Event, time.Time)
}

// WatcherPool control all the watchers
type WatcherPool struct {
    Table    map[string]*alfredWatcher // The Table shows the paths and its according watcher
    List     []*alfredWatcher          // The List includes all the alfredwatchers
    emitter  Emitter
    counter  map[string]int
    watchmen []*Watchman
}

var evManger *WatcherPool

// Initialize a watch pool
func GetManager() *WatcherPool {
    if evManger == nil {
        evManger = &WatcherPool{
            Table:    make(map[string]*alfredWatcher),
            List:     []*alfredWatcher{},
            counter:  make(map[string]int),
            watchmen: []*Watchman{},
        }
    }
    return evManger
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
func (wp *WatcherPool) Attach(man *Watchman, path string) error {
    // Add watchman to pool
    exist := false
    for _, v := range wp.watchmen {
        if v == man {
            exist = true
            break
        }
    }
    if !exist {
        wp.watchmen = append(wp.watchmen, man)
    }

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
        logger.LoggerOf("watchman-logger").Debug("Create new watcher for ", path)
        go func() {
            for {
                ev := <-w.watcher.Event
                if wp.emitter != nil {
                    wp.emitter.Eject(ev, time.Now())
                }
            }
        }()
    }
    if err := w.AddWatch(path); err != nil {
        return err
    } else {
        wp.Table[path] = w
        wp.counter[path] += 1
    }
    logger.LoggerOf("watchman-logger").Debug(path + " is under watching...")
    return nil
}
func (wp *WatcherPool) Dettach(man *Watchman, path string) error {
    if w, ok := wp.Table[path]; !ok {
        return nil
    } else {
        wp.counter[path] -= 1
        if wp.counter[path] == 0 {
            err := w.RemoveWatch(path)
            if err != nil {
                logger.LoggerOf("watchman-logger").Debug(err.Error())
            }
            delete(wp.Table, path)
            delete(wp.counter, path)
            logger.LoggerOf("watchman-logger").Debugf("Remove %v from watching list.", path)
        } else {
            logger.LoggerOf("watchman-logger").Debugf("Remove a reference to %v from watching list.", path)
        }
    }
    return nil
}

func (wp *WatcherPool) shutdown() {
    for fn, _ := range wp.Table {
        wp.Dettach(nil, fn)
    }
    for _, w := range wp.List {
        w.Release()
    }
}
