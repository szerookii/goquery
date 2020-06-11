package goquery

import (
    "errors"
    "time"
    "crypto/rand"
    "bytes"
    "net"
    "strconv"
    "strings"
    "encoding/binary"
)

type QueryStat struct {
    Host string
    Port int
    Motd string
    GameType string
    Version string
    Software string
    Plugins string
    World string
    Online int
    Max int
    Players []string
}

func Query(address string, port int) (*QueryStat, error) {
    
    addrStr := address + ":" + strconv.Itoa(port)
    sessionId, err := GenerateSessionId()
    
    if err != nil {
        return nil, errors.New("Cannot generate session id")
    }
    
    challenge, err := GenerateChallenge(addrStr, sessionId)
    
    if err != nil {
        return nil, errors.New("Cannot generate challenge")
    }
    
    buf := &bytes.Buffer{}
    
    buf.Write([]byte{0xFE, 0xFD, 0x00})
    buf.Write(sessionId)
    binary.Write(buf, binary.BigEndian, challenge)
    buf.Write([]byte{0x00, 0x00, 0x00, 0x00})
    
    bytes := buf.Bytes()
    
    addr, err := net.ResolveUDPAddr("udp4", addrStr)
    if err != nil {
        return nil, err
    }
    
    conn, err := net.DialUDP("udp4", nil, addr)
    if err != nil {
        return nil, err
    }
    
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    defer conn.Close()
    
    _, err = conn.Write(bytes)
    if err != nil {
        return nil, err
    }
    
    rBuf := make([]byte, 1024)
    
    _, err = conn.Read(rBuf)
    if err != nil {
        return nil, err
    }
    
    str := string(rBuf[11:])
    data := strings.Split(str, "\x00\x01player_\x00\x00")
    
    playersArr := strings.Split(data[1], "\u0000")
    players := playersArr[0:len(playersArr)-2]
    
    data = strings.Split(data[0], "\000")
    
    host := data[23]
    port, _ = strconv.Atoi(data[25])
    motd := data[3]
    gametype := strings.ToLower(data[7])
    version := data[9]
    software := data[11]
    plugins := data[13]
    world := data[15]
    online, _ := strconv.Atoi(data[17])
    max, _ := strconv.Atoi(data[19])
    
    query := &QueryStat{}
    query.Host = host
    query.Port = port
    query.Motd = motd
    query.GameType = gametype
    query.Version = version
    query.Software = software
    query.Plugins = plugins
    query.World = world
    query.Online = online
    query.Max = max
    query.Players = players
    
    return query, nil
}

func GenerateChallenge(address string, sessionId []byte) (int32, error) {
    
    buf := &bytes.Buffer{}
    
    buf.Write([]byte{0xFE, 0xFD, 0x09})
    buf.Write(sessionId)
    buf.Write([]byte("")) // why ?
    
    Bytes := buf.Bytes()
    
    addr, err := net.ResolveUDPAddr("udp4", address)
    if err != nil {
        return 0, err
    }
    
    conn, err := net.DialUDP("udp4", nil, addr)
    if err != nil {
        return 0, err
    }
    
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    defer conn.Close()
    
    _, err = conn.Write(Bytes)
    if err != nil {
        return 0, err
    }
    
    rBuf := make([]byte, 1024)
    
    _, err = conn.Read(rBuf)
    if err != nil {
        return 0, err
    }
    
    str := string(rBuf[5:15])
    
    i, err := strconv.Atoi(str)
    if err != nil {
        return 0, err
    }
    
    challenge := int32(i)
    
    return challenge, nil
}

func GenerateSessionId() ([]byte, error) {
    bytes := make([]byte, 4)
    _, err := rand.Read(bytes)
    
    if err != nil {
        return nil, err
    }
    
    return bytes, nil
}