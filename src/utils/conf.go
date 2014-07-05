package utils

import (
    "bitbucket.org/kardianos/osext"
    "goconf.googlecode.com/hg"
    "gopkg.in/yaml.v1"
    "io/ioutil"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

type WatchfileCfg struct {
    Name      string
    Recursive bool
    Expired   int
}
type WatchCfgArr []WatchfileCfg

func (mp WatchCfgArr) Len() int {
    return len(mp)
}
func (mp WatchCfgArr) Swap(i, j int) {
    mp[i], mp[j] = mp[j], mp[i]
}
func (mp WatchCfgArr) Less(i, j int) bool {
    return strings.Count(mp[i].Name, "/") < strings.Count(mp[j].Name, "/")
}

func ConfDir() (string, error) {
    if filename, err := osext.Executable(); err == nil {
        dir := filepath.Dir(filepath.Dir(filename)) + "/conf"
        if _, err = os.Stat(dir); os.IsNotExist(err) {
            return "", err
        } else {
            return dir, nil
        }
    } else {
        return "", err
    }
}

func MainConf() (*conf.ConfigFile, error) {
    if dir, err := ConfDir(); err == nil {
        if cfg, err := conf.ReadConfigFile(dir + "/main.conf"); err == nil {
            return cfg, nil
        } else {
            return nil, err
        }
    } else {
        return nil, err
    }
}

func GetWatchlist() (list []string) {
    if dir, err := ConfDir(); err == nil {
        fwatch := dir + "/watchlist.conf"
        if _, err = os.Stat(fwatch); !os.IsNotExist(err) {
            if data, err := ioutil.ReadFile(fwatch); err == nil {
                var t WatchCfgArr
                if err = yaml.Unmarshal([]byte(data), &t); err == nil {
                    for _, element := range t {
                        level := 1
                        if element.Recursive {
                            level = 64
                        }
                        tlist := Walk(element.Name, level)
                        list = append(list, tlist...)
                    }
                }
            }
        }
    }
    if len(list) == 0 {
        list = Walk("/var/log", 1)
        list = append(list, Walk("/home", 100)...)
    }
    if len(list) > 8000 {
        list = list[0:8000]
    }
    return
}
func GetWhitelist() (list []string) {
    if dir, err := ConfDir(); err == nil {
        fwatch := dir + "/whitelist.conf"
        if _, err = os.Stat(fwatch); !os.IsNotExist(err) {
            if data, err := ioutil.ReadFile(fwatch); err == nil {
                if err = yaml.Unmarshal([]byte(data), &list); err == nil {
                    return
                }
            }
        }
    }
    return
}
func GetExpiredDate(path string) int {
    expired := 30
    if dir, err := ConfDir(); err == nil {
        fwatch := dir + "/watchlist.conf"
        if _, err = os.Stat(fwatch); !os.IsNotExist(err) {
            if data, err := ioutil.ReadFile(fwatch); err == nil {
                var t WatchCfgArr
                if err = yaml.Unmarshal([]byte(data), &t); err == nil {
                    sort.Sort(t)
                    for i := t.Len() - 1; i >= 0; i-- {
                        if nearest := t[i]; strings.HasPrefix(path, nearest.Name) {
                            expired = nearest.Expired
                            break
                        }
                    }
                }
            }
        }
    }
    return expired
}
