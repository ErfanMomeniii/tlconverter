package tlconverter

import (
	"context"
	"net"
)

type Udp struct {
	source net.Addr
}

func (u *Udp) Forward(ctx context.Context, addr net.Addr, pc chan []byte) (err error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		udpAddr, err = net.ResolveUDPAddr(addr.Network(), addr.String())
		if err != nil {
			return
		}
	}

	conn, err := net.DialUDP(udpAddr.Network(), nil, udpAddr)
	if err != nil {
		return
	}

	defer conn.Close()

	go func() {
		defer close(pc)
		for value := range pc {
			n, err := conn.Write(value)
			if n == 0 || err != nil {
				return
			}
		}
	}()

	return nil
}

func (u *Udp) Packets() (error, chan []byte) {
	udpServer, err := net.ListenPacket(u.source.Network(), u.source.String())
	if err != nil {
		return err, nil
	}
	defer udpServer.Close()

	var p chan []byte
	go func() {
		for {
			buf := make([]byte, 1024)
			n, _, err := udpServer.ReadFrom(buf)
			if n == 0 || err != nil {
				continue
			}
			p <- buf
		}
	}()

	return nil, p

}

func (u *Udp) Network() net.Addr {
	return u.source
}
