package kcp

import "net"

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

func DialKCP(nothing *struct{}, closer Closer, incomingData <-chan []byte, outgoingData chan<- []byte, terminationSignal <-chan struct{}, isStream bool) (net.Conn, error) {
	// m := newMultiplexer(closer, config.MaxPacketSize, incomingData, outgoingData, terminationSignal)

	// s := m.newSocket(config, false, !isStream)
	// err := s.startConnect()

	return nil, nil
}

// ListenKCP listens for incoming UDT connections using the existing provided packet connection. It creates a UDT server.
func ListenKCP(nothing *struct{}, closer Closer, incomingData <-chan []byte, outgoingData chan<- []byte, terminationSignal <-chan struct{}) net.Listener {
	// m := newMultiplexer(closer, config.MaxPacketSize, incomingData, outgoingData, terminationSignal)

	// l := &listener{
	// 	m:      m,
	// 	accept: make(chan *udtSocket, 100),
	// 	closed: make(chan struct{}, 1),
	// 	config: config,
	// }

	// m.listenSock = l

	return nil
}
