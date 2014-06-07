package utils

import (
    "os"
    "path/filepath"
    "sort"
    "strings"
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
