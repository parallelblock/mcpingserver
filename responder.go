package main

import (
    "net"
)

type Handshake struct {
    ProtocolVersion uint64
    ServerAddress string
    ServerPort uint16
    SourceIP net.Addr
}

type PlayersEntry struct {
    MaxPlayers int `json:"max"`
    OnlinePlayers int `json:"online"`
    Sample []PlayerEntry `json:"sample,omitempty"`
}

type PlayerEntry struct {
    Name string `json:"name"`
    Uuid string `json:"id"`
}

type VersionEntry struct {
    Name string `json:"name"`
    Protocol uint `json:"protocol"`
}


type PingResponse struct {
    Version VersionEntry `json:"version"`
    Players PlayersEntry `json:"players"`
    Description interface{} `json:"description"`
    Faviconb64 string `json:"favicon,omitempty"`
}

type LegacyPingResponse struct {
    PlayerCount, PlayerMax, ProtocolVersion int
    ServerVersion, Motd string
}

type Responder interface {
    OnConnect(net.Addr) (error)
    RespondPing(*Handshake) (*PingResponse, error)
    RespondJoin(*Handshake, string) (interface{}, error)
    RespondLegacyPing(*Handshake) (*LegacyPingResponse, error)
}

