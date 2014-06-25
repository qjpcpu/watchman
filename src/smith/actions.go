package smith

import (
    "math"
    . "mlog"
    "path/filepath"
    "router"
    "strconv"
    "syscall"
    "time"
    "utils"
)

const (
    MonthLimitAgo = 3
)

var SingleFileDiskOccupyLimit float32

func fromBigFile(clue router.Message) bool {
    limit := 0.001
    if cfg, err := utils.MainConf(); err == nil {
        slimit, _ := cfg.GetString("default", "SingleFileDiskOccupyLimit")
        if f, err := strconv.ParseFloat(slimit, 32); err == nil {
            limit = f
        }
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
    dir := filepath.Dir(clue.FileName)
    timelimit := time.Now().AddDate(0, 0, -30)
    limit := 0.001
    if cfg, err := utils.MainConf(); err == nil {
        slimit, _ := cfg.GetString("default", "TrivialFilesOccupyLimit")
        if f, err := strconv.ParseFloat(slimit, 32); err == nil {
            limit = f
        }
        sd, _ := cfg.GetString("default", "OldFileDateLimit")
        if d, err := strconv.Atoi(sd); err == nil {
            timelimit = time.Now().AddDate(0, 0, -d)
        }
    }
    list, total := utils.Find(dir, 1, 0, timelimit)
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
func erase(files ...string) {
    for _, f := range files {
        Log.Info("Remove", f)
        //os.RemoveAll(f)
    }
}
