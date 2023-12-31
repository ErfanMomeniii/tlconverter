package tlconverter

import (
	"context"
	"net"
)

type Protocol interface {
	Forward(ctx context.Context, addr net.Addr, pc chan []byte) error
	Packets() (error, chan []byte)
	Network() net.Addr
}

type Converter struct {
	source Protocol
	target Protocol
}

func new(source Protocol, target Protocol) *Converter {
	return &Converter{
		source: source,
		target: target,
	}
}

func New(sourceProtocol string, targetProtocol string, sourceAddress string, targetAddress string) *Converter {
	return new(
		protocol(sourceProtocol, sourceAddress),
		protocol(targetProtocol, targetAddress),
	)
}

func (c *Converter) Convert() (err error) {
	go func() {
		err, pc := c.source.Packets()

		err = c.source.Forward(context.Background(), c.target.Network(), pc)
		if err != nil {
			return
		}
	}()

	return
}

func protocol(network string, address string) Protocol {
	switch network {
	case "udp", "udp4", "udp6":
		addr, _ := net.ResolveUDPAddr(network, address)

		return &Udp{addr}
	default:
		addr, _ := net.ResolveTCPAddr(network, address)

		return &Tcp{addr}
	}
}
