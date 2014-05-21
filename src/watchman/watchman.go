package watchman

import (
    "code.google.com/p/go.exp/inotify"
    "errors"
)

type Watchman struct{
    queue_size int
    paths map[string]uint
    queue chan map[string]interface{}
    watcher *inotify.Watcher
}

func NewWatchman(queuesize int) (*Watchman,error){
    iw ,err := inotify.NewWatcher()
    w := &Watchman{
        queue_size:queuesize,
        paths:make(map[string]uint,100),
        watcher:iw,
    }
    return w,err
}
func (man *Watchman) WatchPath(path string,events uint)(error){
    if _,ok:=man.paths[path];ok{
        return errors.New(path+" has already been watched!")
    }
    err := man.watcher.AddWatch(path)
    if err==nil{
        man.paths[path]=events&inotify.IN_ALL_EVENTS
    }
    return err
}
func (man *Watchman) ForgetPath(path string)(error){
    if _,ok:=man.paths[path];!ok{
        return errors.New(path+" has already been watched!")
    }
    err := man.watcher.AddWatch(path)
    if err==nil{
        man.paths[path]=events&inotify.IN_ALL_EVENTS
    }
    return err
}
