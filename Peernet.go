package kcp

import (
	"sync/atomic"
	"time"
)

func PeernetDialWithOptions() {

}

func NewPeerNetUDPSession(conv uint32, dataShards, parityShards int, l *Listener, incomingData <-chan []byte, outgoingData chan<- []byte, terminationSignal <-chan struct{}, block BlockCrypt) *UDPSession {

	sess := new(UDPSession)
	sess.die = make(chan struct{})
	sess.nonce = new(nonceAES128)
	sess.nonce.Init()
	sess.chReadEvent = make(chan struct{}, 1)
	sess.chWriteEvent = make(chan struct{}, 1)
	sess.chSocketReadError = make(chan struct{})
	sess.chSocketWriteError = make(chan struct{})
	sess.l = l
	sess.block = block
	sess.recvbuf = make([]byte, mtuLimit)

	// Custom Peernet channels
	sess.incomingData = incomingData
	sess.outgoingData = outgoingData

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

	sess.kcp = NewKCP(conv, func(buf []byte, size int) {
		if size >= IKCP_OVERHEAD+sess.headerSize {
			sess.output(buf[:size])
		}
	})
	sess.kcp.ReserveBytes(sess.headerSize)

	if sess.l == nil { // it's a client connection
		go sess.readLoop()
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

	return sess
}

// ReadEventPeernet Read event to read input from VirtualPacketConn Channel
func (s *UDPSession) ReadEventPeernet() {

	select {
	// Incoming data
	case <-s.incomingData:
        
	}
}

// WriteEventPeernet Writes event to the send channel from the VirtualPacketConn Channel
func (s *UDPSession) WriteEventPeernet() {

}
