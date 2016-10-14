package main

import (
    "log"
    "net"
)

func CreateSimpleResponder(response *PingResponse, kickResponse interface{}, legacyResponse *LegacyPingResponse) (Responder) {
    return SimpleResponder{response, kickResponse, legacyResponse}
}

type SimpleResponder struct {
    response *PingResponse
    kickResponse interface{}
    legacyResponse *LegacyPingResponse
}

func (s SimpleResponder) OnConnect(a net.Addr) error {
    log.Println("Client connected from", a)
    return nil
}

func (s SimpleResponder) RespondPing(h *Handshake) (*PingResponse, error) {
    log.Println("Responded to ping request from", h.SourceIP)
    return s.response, nil
}

func (s SimpleResponder) RespondJoin(h *Handshake, u string) (interface{}, error) {
    log.Println("Kicking join from", h.SourceIP)
    return s.kickResponse, nil
}

func (s SimpleResponder) RespondLegacyPing(h *Handshake) (*LegacyPingResponse, error) {
    log.Println("Responding to legacy ping request from", h.SourceIP)
    return s.legacyResponse, nil
}
