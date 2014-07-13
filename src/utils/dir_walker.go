package utils

import (
    "os"
    "path/filepath"
    "sort"
    "strings"
    //    "syscall"
    "time"
)

type Path []string

func (p Path) Len() int           { return len(p) }
func (a Path) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (p Path) Less(i, j int) bool { return strings.Count(p[i], "/") < strings.Count(p[j], "/") }

func Walk(root string, level int, exclude ...string) (list []string) {
    visit := func(path string, f os.FileInfo, err error) error {
        if strings.HasPrefix(f.Name(), ".") && f.IsDir() {
            return filepath.SkipDir
        }
        for _, ex := range exclude {
            if strings.HasSuffix(path, ex) && f.IsDir() {
                return filepath.SkipDir
            }
        }
        if f.IsDir() {
            list = append(list, path)
        }
        return nil
    }
    filepath.Walk(root, visit)
    sort.Sort(Path(list))
    end := len(list)
    for i, v := range list {
        if strings.Count(v, "/")-strings.Count(root, "/") > level {
            end = i
            break
        }
    }
    list = list[0:end]
    return
}
func Find(root string, level int, size int64, before time.Time) (list []string, total int64) {
    if _, err := os.Stat(root); os.IsNotExist(err) {
        return
    }
    total = 0
    visit := func(path string, f os.FileInfo, err error) error {
        if level > 0 && strings.Count(path, "/")-strings.Count(root, "/") >= level && f.IsDir() {
            return filepath.SkipDir
        }
        //if t, ok := f.Sys().(*syscall.Stat_t); !f.IsDir() && f.Size() > size && f.ModTime().Before(before) && (ok && time.Unix(t.Atim.Unix()).Before(before)) {
        if !f.IsDir() && f.Size() > size && f.ModTime().Before(before) {
            list = append(list, path)
            total += f.Size()
        }
        return err
    }
    filepath.Walk(root, visit)
    return
}
