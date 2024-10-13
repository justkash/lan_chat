package main

import (
  "fmt"
  "context"
  "bufio"
  "os"
  "flag"
  "bytes"

  "github.com/libp2p/go-libp2p"
  "github.com/libp2p/go-libp2p/core/host"
  "github.com/libp2p/go-libp2p/core/peer"
  "github.com/libp2p/go-libp2p-pubsub"
  "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const (
  ServiceName = "lan-chat"
  Topic = "lan-chat-channel" 
)

type MdnsNotifee struct {
  host host.Host
}

func (notifee *MdnsNotifee) HandlePeerFound(peerAddr peer.AddrInfo) {
  fmt.Printf("Discovered new peer %s\n", peerAddr.ID)
  err := notifee.host.Connect(context.Background(), peerAddr)
  if err != nil {
    fmt.Printf("Error connecting to peer: %s. %s\n", peerAddr.ID, err)
  }
}

func readFromSubscription(ctx context.Context, sub *pubsub.Subscription) {
  for {
    msg, err := sub.Next(ctx)
    if err != nil {
      fmt.Printf("Error reading from subscription: %s\n", err)
      return
    }

    fmt.Printf("%s\n", msg.Data)
  } 
}

func publishMessage(ctx context.Context, topic *pubsub.Topic, name string, id peer.ID, message string) {
  var payloadBuilder bytes.Buffer
  payloadBuilder.WriteString(name)
  payloadBuilder.WriteString("(")
  payloadBuilder.WriteString(id.String())
  payloadBuilder.WriteString("): ")
  payloadBuilder.WriteString(message)

  if err := topic.Publish(ctx, payloadBuilder.Bytes()); err != nil {
    fmt.Printf("Error publishing message: %s\n", err)
  }
}

func joinChat(ctx context.Context, pubsub *pubsub.PubSub, name string, id peer.ID) *pubsub.Topic {
  topic, err := pubsub.Join(Topic)
  if err != nil { 
    fmt.Printf("Error joining topic: %s\n", err)
    return nil
  }

  subscription, err := topic.Subscribe()
  if err != nil {
    fmt.Printf("Error subscribing: %s\n", err)
    return nil
  }

  go readFromSubscription(ctx, subscription)
  return topic
}

func readFromCli(ctx context.Context, topic *pubsub.Topic, name string, id peer.ID) {
  fmt.Print("Type message and enter.\n")
  for {
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    publishMessage(ctx,topic, name, id, text)
  }
}

func main() {
  name := flag.String("name", "gimly", "Nick name to use within the chat")
  flag.Parse()

  // Defaults from https://pkg.go.dev/github.com/libp2p/go-libp2p#New used
	host, err := libp2p.New()
	if err != nil { panic(err) }

  ctx := context.Background()
  pubsubService, err := pubsub.NewGossipSub(ctx, host)
	if err != nil { panic(err) }

  mdnsService := mdns.NewMdnsService(host, ServiceName, &MdnsNotifee{host: host})
  if err = mdnsService.Start(); err != nil { panic(err) }

  selfID := host.ID()
  topic := joinChat(ctx, pubsubService, *name, selfID)
  if topic == nil {
    fmt.Printf("Error joining chat: %s\n")
  }

  fmt.Printf("Name: %s\n", *name)
  fmt.Printf("ID: %s\n", selfID)
  readFromCli(ctx, topic, *name, selfID)
}
