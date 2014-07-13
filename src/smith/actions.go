package smith

import (
    "math"
    "os"
    "path/filepath"
    "router"
    "strconv"
    "strings"
    "syscall"
    "time"
    "utils"
)

const (
    MonthLimitAgo = 3
)

func fromBigFile(clue router.Message) bool {
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

func fromBigDirectory(clue router.Message) ([]string, bool) {
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

func fromWhiteList(clue router.Message) bool {
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

func canErase(files ...string) (yes, no []string) {
    no = []string{}
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
func erase(files ...string) {
    for _, f := range files {
        if strings.HasPrefix(f, "/var/log/") {
            os.Truncate(f, 0)
        } else {
            os.RemoveAll(f)
        }
    }
}
