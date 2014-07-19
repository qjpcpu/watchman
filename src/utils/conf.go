package utils

import (
    "bitbucket.org/kardianos/osext"
    "gopkg.in/yaml.v1"
    "io/ioutil"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

// default configurations
var configurations = struct {
    Main      MainConfig
    WatchList WatchConfigArr
    WhiteList []string
}{}

const ConfigFileName = "watchman.conf"

// main configurations for watchman
type MainConfig struct {
    LogFile                   string
    LogLevel                  string
    Action                    string
    SingleFileDiskOccupyLimit float64
    TrivialFilesOccupyLimit   float64
}

// the watch list
type WatchConfig struct {
    Name      string
    Recursive bool
    Expired   int
}

type WatchConfigArr []WatchConfig

func (mp WatchConfigArr) Len() int {
    return len(mp)
}
func (mp WatchConfigArr) Swap(i, j int) {
    mp[i], mp[j] = mp[j], mp[i]
}
func (mp WatchConfigArr) Less(i, j int) bool {
    return strings.Count(mp[i].Name, "/") < strings.Count(mp[j].Name, "/")
}

func RootDir() (string, error) {
    if filename, err := osext.Executable(); err == nil {
        dir := filepath.Dir(filepath.Dir(filename))
        if _, err = os.Stat(dir); os.IsNotExist(err) {
            return "", err
        } else {
            return dir, nil
        }
    } else {
        return "", err
    }
}

func GetMainConfig() (MainConfig, error) {
    err := LoadConfigurations()
    return configurations.Main, err
}

func GetWatchlist() (list []string) {
    for _, element := range configurations.WatchList {
        level := 1
        if element.Recursive {
            level = 64
        }
        tlist := Walk(element.Name, level)
        list = append(list, tlist...)
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
    return configurations.WhiteList
}
func GetExpiredDate(path string) int {
    expired := 30
    t := configurations.WatchList
    sort.Sort(t)
    for i := t.Len() - 1; i >= 0; i-- {
        if nearest := t[i]; strings.HasPrefix(path, nearest.Name) {
            expired = nearest.Expired
            break
        }
    }
    return expired
}
func LoadConfigurations() error {
    if dir, err := RootDir(); err == nil {
        fwatch := dir + "/conf/" + ConfigFileName
        if _, err = os.Stat(fwatch); !os.IsNotExist(err) {
            if data, err := ioutil.ReadFile(fwatch); err == nil {
                if err = yaml.Unmarshal([]byte(data), &configurations); err == nil {
                    return nil
                } else {
                    return err
                }
            } else {
                return err
            }
        } else {
            return err
        }
    } else {
        return err
    }
}
