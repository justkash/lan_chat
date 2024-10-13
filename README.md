# LAN Chat
This repository holds a simple implementation of a LAN based peer-to-peer chat application that makes use of [libp2p](https://docs.libp2p.io/concepts/introduction/overview/). The application discovers other peers using [mDNS](https://docs.libp2p.io/concepts/discovery-routing/mdns/) and exchanges messages via a [gossip based pubsub](https://docs.libp2p.io/concepts/pubsub/overview/) system.

## Run
To try out the application, you can use the following command to start a single node of the application within the local network. Any additional nodes that are started will automatically connect to the smae chat channel and messages send by a single node will be visible to all other nodes within the LAN.
```go
go run . -name Legolas
```

## Development
This project makes use of [Nix](https://nixos.org/), so it's possible to use the following commands to enter a shell environment with `go` available. The tools available within the Nix shell is defined within the [flake.nix](./flake.nix) file.
```nix
nix develop
```
Furthermore, it's possible to build and run the application using the following commands respectively.
```nix
nix build
nix run
```
