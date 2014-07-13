package smith

import (
    "container/list"
    . "mlog"
    "router"
    "time"
    "utils"
    "watchman"
)

func ScanAbnormal(queue *list.List) {
    c := time.Tick(time.Second)
    action := "info"
    if mcfg, err := utils.GetMainConfig(); err == nil {
        action = mcfg.Action
    }
    for _ = range c {
        erase_list := []string{}
        for {
            if ele := queue.Back(); ele != nil {
                value := queue.Remove(ele)
                msg := value.(router.Message)
                if fromWhiteList(msg) {
                    Log.Infof("%v is on  white list,pass.", msg.FileName)
                } else if fromBigFile(msg) {
                    if canEraseInstant(msg.FileName) {
                        erase_list = append(erase_list, msg.FileName)
                    } else {
                        Log.Warningf("Big file found and I dare not del.(%v:%v)", msg.FileName, msg.Size)
                    }
                } else if can_del, ok := fromBigDirectory(msg); ok {
                    if yes, _ := canErase(can_del...); len(yes) > 0 {
                        erase_list = append(erase_list, yes...)
                    }
                } else {
                    Log.Debugf("You're good %s(%s), let you go.", msg.FileName, watchman.HumanReadable(msg.Event))
                }
            } else {
                break
            }
        }
        if len(erase_list) > 0 {
            Log.Infof("[%s] Remove %v", action, erase_list)
            switch action {
            case "info":
            case "remove":
                erase(erase_list...)
            }
        }
    }
}
