package watchman

import (
    "alfred"
    "os"
    "testing"
)

func TestWatchman(t *testing.T) {
    alfred.Boot()
    man, err := NewWatchman()
    if err != nil {
        t.Fatal(err)
    }
    if err = man.WatchPath("/tmp", IN_ALL_EVENTS); err != nil {
        t.Fatal(err)
    }
    if sm, err := man.PullEvent(); err.Error() != "SYSTEM" || sm.FileName != "SUCCESS:+/tmp" {
        t.Fatal("Watch /tmp fail")
    }
    fn := "/tmp/create_for_watchman_test"
    f, _ := os.Create(fn)
    f.Close()
    os.Remove(fn)
    fopen, fremove := false, false
    for {
        if m, err := man.PullEvent(); err == nil {
            if m.Event&IN_OPEN != 0 {
                fopen = true
            }
            if m.Event&IN_DELETE != 0 {
                fremove = true
            }
        }
        if fopen && fremove {
            break
        }
    }
    if err = man.ForgetPath("/tmp"); err != nil {
        t.Fatal(err)
    }
    man.Release()
    alfred.Shutdown()
}
