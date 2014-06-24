package main

import (
    "gopkg.in/yaml.v1"
    "io/ioutil"
    "os"
    "utils"
)

type onwatch struct {
    Name      string
    Recursive bool
}

func getWatchlist() (list []string) {
    if dir, err := utils.ConfDir(); err == nil {
        fwatch := dir + "/watchlist.conf"
        if _, err = os.Stat(fwatch); !os.IsNotExist(err) {
            if data, err := ioutil.ReadFile(fwatch); err == nil {
                var t []onwatch
                if err = yaml.Unmarshal([]byte(data), &t); err == nil {
                    for _, element := range t {
                        level := 1
                        if element.Recursive {
                            level = 64
                        }
                        tlist := utils.Walk(element.Name, level)
                        list = append(list, tlist...)
                    }
                }
            }
        }
    }
    if len(list) == 0 {
        list = utils.Walk("/var/log", 1)
        list = append(list, utils.Walk("/home", 100)...)
    }
    if len(list) > 8000 {
        list = list[0:8000]
    }
    return
}
