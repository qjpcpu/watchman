package router

import (
    "errors"
    "gopkg.in/redis.v1"
)

const SYS_ID = "_alfred_"

type RouterCli struct {
    id     string
    client *redis.Client
    pubsub *redis.PubSub
}

type BuildClient func() *redis.Client

func DefaultBuildClient() *redis.Client {
    return redis.NewTCPClient(&redis.Options{
        Addr: ":6379",
    })
}

func NewRouterCli(uid string, builder BuildClient) *RouterCli {
    cli := &RouterCli{
        id:     uid,
        client: builder(),
    }
    cli.pubsub = cli.client.PubSub()
    return cli
}

func (rc *RouterCli) Write(to, msg string) error {
    pub := rc.client.Publish(to, msg)
    return pub.Err()
}
func (rc *RouterCli) Subscribe(path string) {
    rc.pubsub.Subscribe(path)
}
func (rc *RouterCli) Unsubscribe(path string) {
    rc.pubsub.Unsubscribe(path)
}
func (rc *RouterCli) Read() (string, error) {
    msg, err := rc.pubsub.Receive()
    if err != nil {
        return "", err
    }
    payload, ok := msg.(*redis.Message)
    if !ok {
        return "", errors.New("Not a message.")
    }
    return payload.Payload, nil
}

func (rc *RouterCli) Close() error {
    rc.pubsub.Close()
    return rc.client.Close()
}
