package main

import (
    "bufio"
    "bytes"
    "encoding/binary"
    "encoding/json"
    "io"
    "net"
)

func readString(b *bufio.Reader) (string, error) {
    strlen, err := binary.ReadUvarint(b)
    if err != nil {
        return "", err
    }
    strRaw := make([]uint8, strlen, strlen)
    _, err = b.Peek(int(strlen))
    err = binary.Read(b, binary.BigEndian, strRaw)
    return string(strRaw), err
}

 type PlayerConn struct {
    pingServer *PingServer

    inBuf *bufio.Reader
    conn net.Conn

    errorHandler PingServerErrorHandler
}

func (pc *PlayerConn) isLegacy() (bool, error) {
    headerPacket, err := pc.inBuf.Peek(1)
    if err != nil {
        return false, err
    }
    return headerPacket[0] & 0xFE == 0xFE, nil // 0xFE or 0xFF
}

func (pc *PlayerConn) handleConnectionWThrow() error {
    defer pc.conn.Close()

    err := pc.pingServer.Responder.OnConnect(pc.conn.RemoteAddr())
    if err != nil {
        return err
    }

    lgcy, err := pc.isLegacy()
    if err != nil {
        return err
    }

    if lgcy {
        return pc.handleLegacy()
    } else {
        return pc.handleModern()
    }
}

func (pc *PlayerConn) handleConnection() {
    err := pc.handleConnectionWThrow()
    if err != nil {
        pc.errorHandler(err)
    }
}

func (p *PlayerConn) handleModern() error {
    packet, err := p.readPacket()
    if err != nil {
        return err
    }

    packetId, err := binary.ReadUvarint(packet)
    if err != nil || packetId != 0 {
        return err
    }

    protoVersion, err := binary.ReadUvarint(packet)
    if err != nil {
        return err
    }

    hostname, err := readString(packet)
    if err != nil {
        return err
    }

    var port uint16
    err = binary.Read(packet, binary.BigEndian, &port)
    if err != nil {
        return err
    }

    nextState, err := binary.ReadUvarint(packet)
    if err != nil {
        return err
    }

    hs := Handshake{protoVersion, hostname, port, p.conn.RemoteAddr()}

    if nextState == 1 {
        // status
        return p.handlePing(&hs)
    } else if nextState == 2 {
        // wait for login
        return p.handleJoin(&hs)
    } else {
        return nil
    }
}


func (p *PlayerConn) handleJoin(h *Handshake) error {
    packet, err := p.readPacket()
    if err != nil {
        return err
    }
    packetId, err := binary.ReadUvarint(packet)
    if err != nil || packetId != 0 {
        return err
    }

    username, err := readString(packet)
    if err != nil {
        return err
    }

    resp, err := p.pingServer.Responder.RespondJoin(h, username)
    if err != nil {
        return err
    }
    return p.writeJsonPacket(0, resp)
}

func (p *PlayerConn) handlePing(h *Handshake) error {
    packet, err := p.readPacket()
    packetId, err := binary.ReadUvarint(packet)

    if err != nil || packetId != 0 {
        return err
    }

    resp, err := p.pingServer.Responder.RespondPing(h)
    if err != nil {
        return err
    }

    err = p.writeJsonPacket(0, resp)

    packet, err = p.readPacket()
    if err != nil {
        return err
    }
    packetId, err = binary.ReadUvarint(packet)

    if err != nil || packetId != 1 {
        return err
    }

    var secretNumberBuf bytes.Buffer
    i, err := secretNumberBuf.ReadFrom(packet)
    if err != nil || i > 8 {
        return err
    }

    return p.writePacket(1, &secretNumberBuf)
}

func (p *PlayerConn) readPacket() (*bufio.Reader, error) {
    packlen, err := binary.ReadUvarint(p.inBuf)
    if err != nil {
        return nil, err
    }

    if packlen > 256 || packlen < 0 {
        return nil, io.ErrShortBuffer
    }

    return bufio.NewReader(io.LimitReader(p.inBuf, int64(packlen))), nil
}

func (p *PlayerConn) writeJsonPacket(packetid uint64, data interface{}) error {
    marshalledData, err := json.Marshal(data)
    if err != nil {
        return err
    }

    strLen := make([]byte, 9, 9)
    strLenLen := binary.PutUvarint(strLen, uint64(len(marshalledData)))
    var buf bytes.Buffer
    buf.Write(strLen[:strLenLen])
    buf.Write(marshalledData)
    return p.writePacket(packetid, &buf)
}

func (p *PlayerConn) writePacket(packetid uint64, data *bytes.Buffer) error {
    idBuf := make([]byte, 9, 9)
    idBufLen := binary.PutUvarint(idBuf, packetid)
    lenBuf := make([]byte, 9, 9)
    lenBufLen := binary.PutUvarint(lenBuf, uint64(idBufLen + data.Len()))
    err := p.writeSlices(lenBuf[:lenBufLen], idBuf[:idBufLen])
    if err == nil {
        _, err = data.WriteTo(p.conn)
    }
    return err
}

func (p *PlayerConn) writeSlices(slices ...[]byte) (err error) {
    for i := 0; err == nil && i < len(slices); i++ {
        _, err = p.conn.Write(slices[i])
    }
    return
}
