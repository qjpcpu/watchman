package alfred

var pool *WatcherPool

func init() {
    pool = initPool()
    pool.boot()
}
