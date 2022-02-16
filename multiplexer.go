package kcp

import "time"

// Config controls behavior of sockets created with it
type Config struct {
	MaxPacketSize  uint          // Upper limit on maximum packet size (0 = unlimited)
	MaxFlowWinSize uint          // maximum number of unacknowledged packets to permit (minimum 32)
	SynTime        time.Duration // SynTime

	// CanAccept           func(hsPacket *packet.HandshakePacket) error // can this listener accept this connection?
	// CongestionForSocket func(sock *udtSocket) CongestionControl      // create or otherwise return the CongestionControl for this socket
}

// A multiplexer is a single UDT socket over a single PacketConn.
type multiplexer struct {
	//socketID          uint32          // Socket ID
	maxPacketSize     uint            // the Maximum Transmission Unit of packets sent from this address
	incomingData      <-chan []byte   // source to read packets from
	outgoingData      chan<- []byte   // destination to send packets to
	terminationSignal <-chan struct{} // external termination signal to watch
	closer            Closer          // external closer to call in case the local socket/listener closes
}

// The closer is called when the socket/listener closes. The terminationSignal is an external (upstream) signal to watch for.
func newMultiplexer(closer Closer, maxPacketSize uint, incomingData <-chan []byte, outgoingData chan<- []byte, terminationSignal <-chan struct{}) (m *multiplexer) {
	m = &multiplexer{
		maxPacketSize:     maxPacketSize,
		closer:            closer,
		incomingData:      incomingData,
		outgoingData:      outgoingData,
		terminationSignal: terminationSignal,
	}

	return
}
