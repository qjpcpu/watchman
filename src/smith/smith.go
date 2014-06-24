package smith

import (
    "container/list"
    . "mlog"
    "router"
    "time"
    "watchman"
)

func ScanAbnormal(queue *list.List) {
    c := time.Tick(1 * time.Second)
    for _ = range c {
        erase_list := []string{}
        for {
            if ele := queue.Back(); ele != nil {
                value := queue.Remove(ele)
                msg := value.(router.Message)
                if fromWhiteList(msg) {
                    Log.Debugf("%v is on  white list,pass.", msg.FileName)
                } else if fromBigFile(msg) {
                    if yes, _ := canErase(msg.FileName); len(yes) > 0 {
                        Log.Debugf("I will&can kill %s(%s)", msg.FileName, watchman.HumanReadable(msg.Event))
                        erase_list = append(erase_list, msg.FileName)
                    }
                } else if can_del, ok := fromBigDirectory(msg); ok {
                    if yes, _ := canErase(can_del...); len(yes) > 0 {
                        Log.Debugf("I will&can kill %v(%s)", yes, watchman.HumanReadable(msg.Event))
                        erase_list = append(erase_list, yes...)
                    }
                } else {
                    Log.Debugf("You're good %s(%s), let you go.", msg.FileName, watchman.HumanReadable(msg.Event))
                }
            } else {
                break
            }
        }
        erase(erase_list...)
    }
}
