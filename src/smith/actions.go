package smith

import (
    "math"
    . "mlog"
    "router"
    "syscall"
)

const (
    SizeLimitPercentage = 0.001
    MonthLimitAgo       = 3
)

func fromBigFile(clue router.Message) bool {
    file := clue.FileName
    fs := syscall.Statfs_t{}
    syscall.Statfs(file, &fs)
    if percentage := float64(clue.Size) / (float64(fs.Bsize) * float64(fs.Blocks)); !math.IsInf(percentage, 1) && percentage > SizeLimitPercentage {
        return true
    }
    return false
}

func fromBigDirectory(clue router.Message) ([]string, bool) {
    return []string{}, false
}

func fromWhiteList(clue router.Message) bool {
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
    }
}
