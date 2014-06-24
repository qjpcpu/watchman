package utils

import (
    "bitbucket.org/kardianos/osext"
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
