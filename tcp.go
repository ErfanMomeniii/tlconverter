package tlconverter

import (
	"context"
	"net"
)

type Tcp struct {
	source net.Addr
}

func (t *Tcp) Forward(ctx context.Context, addr net.Addr, pc chan []byte) (err error) {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		tcpAddr, err = net.ResolveTCPAddr(addr.Network(), addr.String())
		if err != nil {
			return
		}
	}

	conn, err := net.DialTCP(tcpAddr.Network(), nil, tcpAddr)
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

func (t *Tcp) Packets() (error, chan []byte) {
	tcpServer, err := net.ListenPacket(t.source.Network(), t.source.String())
	if err != nil {
		return err, nil
	}
	defer tcpServer.Close()

	var p chan []byte
	go func() {
		for {
			buf := make([]byte, 1024)
			n, _, err := tcpServer.ReadFrom(buf)
			if n == 0 || err != nil {
				continue
			}
			p <- buf
		}
	}()

	return nil, p
}

func (t *Tcp) Network() net.Addr {
	return t.source
}
