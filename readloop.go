package kcp

func (s *UDPSession) readLoop() {
	for {
		var buf []byte
		select {
		case buf = <-s.m.incomingData:
			s.packetInput(buf)
		case <-s.m.terminationSignal:
			s.Close()
			break
		case <-s.die:
			break
		}
	}
}

func (l *Listener) monitor() {
	for {
		var buf []byte
		select {
		case buf = <-l.m.incomingData:
			l.packetInput(buf)
		case <-l.m.terminationSignal:
			l.Close()
			break
		case <-l.die:
			break
		}
	}
}
