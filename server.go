package main

import (
    "bufio"
    "net"
    "time"
)

var sectionChar = "§"
type PingServer struct {
    bindAddr string
    Responder Responder // only responderhook can be hotswapped
    bindConn net.Listener
}

func CreatePingServer(bindAddr string, hook Responder) *PingServer {
    return &PingServer{bindAddr, hook, nil}
}

func (ps *PingServer) Bind() (err error) {
    ps.bindConn, err = net.Listen("tcp", ps.bindAddr)
    return
}

type PingServerErrorHandler func(error)

func (ps *PingServer) AcceptConnections(handler PingServerErrorHandler) (err error) {
    for err == nil {
        err = ps.AcceptConnection(handler)
    }
    return
}

func (ps *PingServer) AcceptConnection(handler PingServerErrorHandler) error {
    conn, err := ps.bindConn.Accept()
    if err != nil {
        return err
    }
    conn.SetDeadline(time.Now().Add(5 * time.Second))

    connInBuffer := bufio.NewReader(conn)

    playerConn := PlayerConn{ps, connInBuffer, conn, handler}
    go playerConn.handleConnection()
    return nil
}
