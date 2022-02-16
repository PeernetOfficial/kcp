package kcp

import (
	"sync/atomic"

	"golang.org/x/net/ipv4"
)

func (s *UDPSession) defaultTx(txqueue []ipv4.Message) {
	nbytes := 0
	npkts := 0
	for k := range txqueue {
		s.m.outgoingData <- txqueue[k].Buffers[0]
		nbytes += len(txqueue[k].Buffers[0])
		npkts++
	}
	atomic.AddUint64(&DefaultSnmp.OutPkts, uint64(npkts))
	atomic.AddUint64(&DefaultSnmp.OutBytes, uint64(nbytes))
}

func (s *UDPSession) tx(txqueue []ipv4.Message) {
	s.defaultTx(txqueue)
}
