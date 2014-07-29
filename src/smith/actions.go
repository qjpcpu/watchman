package smith

import (
    "alfred"
    "bytes"
    "encoding/json"
    "fmt"
    "math"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "syscall"
    "time"
    "utils"
)

const (
    MonthLimitAgo = 3
)

func fromBigFile(clue alfred.Message) bool {
    limit := 0.001
    if cfg, err := utils.GetMainConfig(); err == nil {
        limit = cfg.SingleFileDiskOccupyLimit
    }
    file := clue.FileName
    fs := syscall.Statfs_t{}
    syscall.Statfs(file, &fs)
    if percentage := float64(clue.Size) / (float64(fs.Bsize) * float64(fs.Blocks)); !math.IsInf(percentage, 1) && percentage > limit {
        return true
    }
    return false
}

func fromBigDirectory(clue alfred.Message) ([]string, bool) {
    level := 1
    dir := filepath.Dir(clue.FileName)
    if _, err := strconv.Atoi(filepath.Base(dir)); err == nil {
        dir = filepath.Dir(dir)
        level = 2
    }
    tl := utils.GetExpiredDate(dir)
    timelimit := time.Now().AddDate(0, 0, -tl)
    limit := 0.001
    if cfg, err := utils.GetMainConfig(); err == nil {
        limit = cfg.TrivialFilesOccupyLimit
    }
    list, total := utils.Find(dir, level, 1000000, timelimit)
    fs := syscall.Statfs_t{}
    syscall.Statfs(dir, &fs)
    if percentage := float64(total) / (float64(fs.Bsize) * float64(fs.Blocks)); !math.IsInf(percentage, 1) && percentage > limit {
        return list, true
    }
    return []string{}, false
}

func fromWhiteList(clue alfred.Message) bool {
    list := utils.GetWhitelist()
    if list != nil {
        for _, e := range list {
            if e == clue.FileName {
                return true
            }
        }
    }
    return false
}

func canErase(files ...alfred.Message) (yes, no []alfred.Message) {
    no = []alfred.Message{}
    yes = files
    return
}
func canEraseInstant(file string) bool {
    del := false
    if strings.HasPrefix(file, "/var/log/") {
        del = true
    }
    return del
}
func erase(files ...alfred.Message) {
    for _, mf := range files {
        f := mf.FileName
        if strings.HasPrefix(f, "/var/log/") {
            os.Truncate(f, 0)
        } else {
            os.RemoveAll(f)
        }
    }
}
func printState(files ...alfred.Message) {
    dir, err := utils.RootDir()
    if err != nil {
        return
    }
    fi, err := os.Create(dir + "/status/watchman.json")
    if err != nil {
        return
    }
    defer fi.Close()
    type Kill struct {
        Name string
        Size string
    }
    kills := []Kill{}
    for _, f := range files {
        k := Kill{
            Name: f.FileName,
            Size: fmt.Sprintf("%vM", f.Size/(1024*1024)),
        }
        kills = append(kills, k)
    }
    b, err := json.Marshal(kills)
    if err != nil {
        return
    }
    var out bytes.Buffer
    json.Indent(&out, b, "", "\t")
    out.WriteTo(fi)
}
