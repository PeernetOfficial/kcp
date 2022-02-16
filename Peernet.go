package kcp

import (
	"crypto/rand"
	"encoding/binary"
	"net"
)

// Closer provides a status code indicating why the closing happens.
type Closer interface {
	Close(reason int) error       // Close is called when the socket is actually closed.
	CloseLinger(reason int) error // CloseLinger is called when the socket indicates to be closed soon, after the linger time.
}

// The termination reason is passed on to the close function
const (
	TerminateReasonListenerClosed     = 1000 // Listener: The listener.Close function was called.
	TerminateReasonLingerTimerExpired = 1001 // Socket: The linger timer expired. Use CloseLinger to know the actual closing reason.
	TerminateReasonConnectTimeout     = 1002 // Socket: The connection timed out when sending the initial handshake.
	TerminateReasonRemoteSentShutdown = 1003 // Remote peer sent a shutdown message.
	TerminateReasonSocketClosed       = 1004 // Send: Socket closed. Called udtSocket.Close().
	TerminateReasonInvalidPacketIDAck = 1005 // Send: Invalid packet ID received in ACK message.
	TerminateReasonInvalidPacketIDNak = 1006 // Send: Invalid packet ID received in NAK message.
	TerminateReasonCorruptPacketNak   = 1007 // Send: Invalid NAK packet received.
	TerminateReasonSignal             = 1008 // Send: Terminate signal. Called udtSocket.Terminate().
)

const dataShardsDefault = 10
const parityShardsDefault = 3

func DialKCP(config *Config, closer Closer, incomingData <-chan []byte, outgoingData chan<- []byte, terminationSignal <-chan struct{}, isStream bool) (net.Conn, error) {
	m := newMultiplexer(closer, config.MaxPacketSize, incomingData, outgoingData, terminationSignal)

	var block BlockCrypt
	block = nil // No encryption for now.
	dataShards := dataShardsDefault
	parityShards := parityShardsDefault
	var convid uint32
	binary.Read(rand.Reader, binary.LittleEndian, &convid)

	s := newUDPSession(m, convid, dataShards, parityShards, nil, block)

	return s, nil
}

// ListenKCP listens for incoming KCP connections using the existing provided packet connection. It creates a KCP server.
func ListenKCP(config *Config, closer Closer, incomingData <-chan []byte, outgoingData chan<- []byte, terminationSignal <-chan struct{}) net.Listener {
	m := newMultiplexer(closer, config.MaxPacketSize, incomingData, outgoingData, terminationSignal)

	l := &Listener{
		m: m,
	}

	l.chAccepts = make(chan *UDPSession, acceptBacklog)
	l.die = make(chan struct{})
	l.dataShards = dataShardsDefault
	l.parityShards = parityShardsDefault
	l.block = nil // No encryption for now.
	go l.monitor()

	return l
}
