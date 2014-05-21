package alfred

import (
	"code.google.com/p/go.exp/inotify"
	"errors"
	"log"
	"reflect"
)

const (
	MAX_PATH_PER_WATCHER = 87
	MAX_WATCHER          = 94
)

type AlfredWatcher struct {
	watcher *inotify.Watcher
	list    map[string]uint32
}

func NewAlfredWatcher() *AlfredWatcher {
	w, _ := inotify.NewWatcher()
	aw := &AlfredWatcher{
		watcher: w,
		list:    make(map[string]uint32),
	}
	return aw
}
func (aw *AlfredWatcher) AddWatch(path string) error {
	if aw.Size() >= MAX_PATH_PER_WATCHER {
		return errors.New("Watch path full for this watcher.")
	}
	err := aw.watcher.Watch(path)
	if err == nil {
		aw.list[path] = inotify.IN_ALL_EVENTS
	}
	return err
}
func (aw *AlfredWatcher) RemoveWatch(path string) error {
	if _, ok := aw.list[path]; ok {
		delete(aw.list, path)
	}
	if err := aw.watcher.RemoveWatch(path); err != nil {
		return err
	}
	return nil
}
func (aw *AlfredWatcher) Release() error {
	return aw.watcher.Close()
}

func (aw *AlfredWatcher) Size() int {
	return len(aw.list)
}

// WatcherPool control all the watchers
type WatcherPool struct {
	Table  map[string]*AlfredWatcher
	List   []*AlfredWatcher
	Signal chan map[string]string
}

func (wp *WatcherPool) FileList() []string {
	list := make([]string, len(wp.Table))
	i := 0
	for k, _ := range wp.Table {
		list[i] = k
		i += 1
	}
	return list
}
func (wp *WatcherPool) Attach(path string) error {
	if _, ok := wp.Table["path"]; ok {
		return nil
	}
	var w *AlfredWatcher
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
		w = NewAlfredWatcher()
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

func initPool() *WatcherPool {
	return &WatcherPool{
		make(map[string]*AlfredWatcher),
		[]*AlfredWatcher{},
		make(chan map[string]string),
	}
}

func (wp *WatcherPool) GetDefaultPaths() []string {
	return []string{"/tmp", "/home/jason/", "/home/jason/tmp/"}
}
func (wp *WatcherPool) PullEvent() *inotify.Event {
	cases := make([]reflect.SelectCase, len(wp.List)+1)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(wp.Signal)}
	for i, ch := range wp.List {
		cases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.watcher.Event)}
	}
	for {
		chosen, value, _ := reflect.Select(cases)
		if chosen == 0 {
            msg:=value.Interface().(map[string]string)
            if path:=msg["PATH"];msg["ACTION"]=="ADD"{

            }else if msg["ACTION"]=="REMOVE"{

            }
		} else {
			return value.Interface().(*inotify.Event)
		}
	}
}
func (wp *WatcherPool) Boot() {
	for _, fn := range wp.GetDefaultPaths() {
		if err := wp.Attach(fn); err != nil {
			log.Println(err)
		}
	}
	go func() {
		for {
			ev := wp.PullEvent()
			log.Println(ev)
		}
	}()
}

var Pool *WatcherPool

func init() {
	Pool = initPool()
	Pool.Boot()
}
