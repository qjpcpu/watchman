package alfred

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
)

// In general, the max watcher object is 8192, so MAX_WATCHER * MAX_PATH_PER_WATCHER < 8192
const (
    MAX_PATH_PER_WATCHER = 87
    MAX_WATCHER          = 94
)

// Alfredwatcher wrapps inotify's wachter, and add the listen file list.
type alfredWatcher struct {
    watcher *inotify.Watcher
    list    map[string]uint32
}

// Create new alfredwatcher, all the methods of alfredwatcher should been invoked by alfred, not by client directly.
func newAlfredWatcher() *alfredWatcher {
    w, _ := inotify.NewWatcher()
    aw := &alfredWatcher{
        watcher: w,
        list:    make(map[string]uint32),
    }
    return aw
}

// Add a path to listen list
func (aw *alfredWatcher) AddWatch(path string) error {
    if aw.Size() >= MAX_PATH_PER_WATCHER {
        return errors.New("Watch path full for this watcher.")
    }
    err := aw.watcher.Watch(path)
    if err == nil {
        aw.list[path] = inotify.IN_ALL_EVENTS
    }
    return err
}

// Remove path from listen list
func (aw *alfredWatcher) RemoveWatch(path string) error {
    if _, ok := aw.list[path]; ok {
        delete(aw.list, path)
    }
    if err := aw.watcher.RemoveWatch(path); err != nil {
        return err
    }
    return nil
}

// Close an alfredwatcher
func (aw *alfredWatcher) Release() error {
    return aw.watcher.Close()
}

// Return the length of listen list
func (aw *alfredWatcher) Size() int {
    return len(aw.list)
}
