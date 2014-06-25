package utils

import (
    "bitbucket.org/kardianos/osext"
    "goconf.googlecode.com/hg"
    "gopkg.in/yaml.v1"
    "io/ioutil"
    "os"
    "path/filepath"
)

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
                var t []struct {
                    Name      string
                    Recursive bool
                }
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
