package kcp

import (
	"crypto/rand"
	"encoding/binary"
	"net"
	"sync/atomic"
	"time"
)

// Closer provides a status code indicating why the closing happens.
type Closer interface {
	Close(reason int) error       // Close is called when the socket is actually closed.
	CloseLinger(reason int) error // CloseLinger is called when the socket indicates to be closed soon, after the linger time.
}

func PeerNetClientUDPSession(dataShards, parityShards int, closer Closer, incomingData <-chan []byte, outgoingData chan<- []byte, terminationSignal <-chan struct{}, block BlockCrypt) (net.Conn, error) {

	sess := new(UDPSession)
	sess.die = make(chan struct{})
	sess.nonce = new(nonceAES128)
	sess.nonce.Init()
	sess.chReadEvent = make(chan struct{}, 1)
	sess.chWriteEvent = make(chan struct{}, 1)
	sess.chSocketReadError = make(chan struct{})
	sess.chSocketWriteError = make(chan struct{})
	sess.l = nil
	sess.block = block
	sess.recvbuf = make([]byte, mtuLimit)

	// Custom Peernet channels
	sess.incomingData = incomingData
	sess.outgoingData = outgoingData
	sess.terminationSignal = terminationSignal

	// Custom peernet closer
	sess.closer = closer

	// FEC codec initialization
	sess.fecDecoder = newFECDecoder(dataShards, parityShards)
	if sess.block != nil {
		sess.fecEncoder = newFECEncoder(dataShards, parityShards, cryptHeaderSize)
	} else {
		sess.fecEncoder = newFECEncoder(dataShards, parityShards, 0)
	}

	// calculate additional header size introduced by FEC and encryption
	if sess.block != nil {
		sess.headerSize += cryptHeaderSize
	}
	if sess.fecEncoder != nil {
		sess.headerSize += fecHeaderSizePlus2
	}

	var conv uint32
	binary.Read(rand.Reader, binary.LittleEndian, &conv)

	sess.kcp = NewKCP(conv, func(buf []byte, size int) {
		if size >= IKCP_OVERHEAD+sess.headerSize {
			sess.output(buf[:size])
		}
	})
	sess.kcp.ReserveBytes(sess.headerSize)

	if sess.l == nil { // it's a client connection
		//go sess.readLoop()
		// Read event peernet
		go sess.ReadEventPeernet()
		atomic.AddUint64(&DefaultSnmp.ActiveOpens, 1)
	} else {
		atomic.AddUint64(&DefaultSnmp.PassiveOpens, 1)
	}

	// start per-session updater
	SystemTimedSched.Put(sess.update, time.Now())

	currestab := atomic.AddUint64(&DefaultSnmp.CurrEstab, 1)
	maxconn := atomic.LoadUint64(&DefaultSnmp.MaxConn)
	if currestab > maxconn {
		atomic.CompareAndSwapUint64(&DefaultSnmp.MaxConn, maxconn, currestab)
	}

	return sess, nil
}

func PeerNetServerUDPSession(closer Closer, incomingData <-chan []byte, outgoingData chan<- []byte, terminationSignal <-chan struct{}, block BlockCrypt, dataShards, parityShards int) net.Listener {
	l := new(Listener)
	//l.conn = conn
	//l.ownConn = ownConn
	l.sessions = make(map[string]*UDPSession)
	l.chAccepts = make(chan *UDPSession, acceptBacklog)
	l.chSessionClosed = make(chan net.Addr)
	l.die = make(chan struct{})
	l.dataShards = dataShards
	l.parityShards = parityShards
	l.block = block
	l.chSocketReadError = make(chan struct{})
	l.closer = closer
	l.outgoingData = outgoingData
	l.incomingData = incomingData
	l.terminationSignal = terminationSignal

	return l
}

// ReadEventPeernet Read event to read input from VirtualPacketConn Channel
func (s *UDPSession) ReadEventPeernet() {
	for {
		var buf []byte
		select {
		// Incoming data
		case buf = <-s.incomingData:
			s.packetInput(buf)
		case <-s.terminationSignal:
			return
		}

		s.Read(buf)
	}
}
