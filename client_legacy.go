package main

import (
    "encoding/binary"
    "fmt"
    "unicode/utf16"
)

func (pc *PlayerConn) handleLegacy() error {
    header, err := pc.inBuf.ReadByte()

    if err != nil || header != 0xfe {
        return err
    }

    // weed out 1.4/1.5
    // isn't overly accurate, but its better than waiting for IO to timeout
    buffered := pc.inBuf.Buffered()
    if buffered == 0 {
        return pc.handle14()
    } else if buffered == 1 {
        return pc.handle15()
    }

    _, err = pc.inBuf.Discard(28)
    if err != nil {
        return err
    }

    var protVersion uint8
    err = binary.Read(pc.inBuf, binary.BigEndian, &protVersion)
    if err != nil {
        return err
    }

    var hostnameLen int16
    err = binary.Read(pc.inBuf, binary.BigEndian, &hostnameLen)
    if err != nil {
        return err
    }

    hostnameLen -= 7

    hostname := make([]uint16, hostnameLen, hostnameLen)
    err = binary.Read(pc.inBuf, binary.BigEndian, &hostname)
    if err != nil {
        return err
    }
    hostnameStr := string(utf16.Decode(hostname))

    var port uint32
    err = binary.Read(pc.inBuf, binary.BigEndian, &port)
    if err != nil {
        return err
    }

    handshake := Handshake{uint64(protVersion), hostnameStr, uint16(port), pc.conn.RemoteAddr()}
    response, err := pc.pingServer.Responder.RespondLegacyPing(&handshake)
    if err != nil {
        return err
    }
    return pc.writeLegacyPing(response)
}

func (pc *PlayerConn) handle15() error {
    print("15")
    handshake := Handshake{0, "", 0, pc.conn.RemoteAddr()}
    resp, err := pc.pingServer.Responder.RespondLegacyPing(&handshake)
    if err != nil {
        return err
    }

    return pc.writeLegacyPing(resp)
}

func (pc *PlayerConn) handle14() error {
    print("14")
    handshake := Handshake{0, "", 0, pc.conn.RemoteAddr()}
    resp, err := pc.pingServer.Responder.RespondLegacyPing(&handshake)
    if err != nil {
        return err
    }

    return pc.writeLegacy(fmt.Sprintf("%s" + sectionChar + "%d" + sectionChar + "%d",
        resp.Motd, resp.PlayerCount, resp.PlayerMax))
}

func (pc *PlayerConn) writeLegacyPing(r *LegacyPingResponse) error {
    return pc.writeLegacy(fmt.Sprintf(sectionChar + "1\u0000%d\u0000%s\u0000%s\u0000%d\u0000%d",
        r.ProtocolVersion, r.ServerVersion, r.Motd, r.PlayerCount, r.PlayerMax))
}

func (pc *PlayerConn) writeLegacy(data string) error {
    err := binary.Write(pc.conn, binary.BigEndian, uint8(255))
    if err != nil {
        return err
    }

    err = binary.Write(pc.conn, binary.BigEndian, uint16(len(data) - 1))
    if err != nil {
        return err
    }
    return binary.Write(pc.conn, binary.BigEndian, utf16.Encode([]rune(data)))
}
