package tlconverter

import (
	"context"
	"net"
	"sync"
)

type Protocol interface {
	Forward(ctx context.Context, addr net.Addr, pc chan<- []byte) error
	Packets() chan<- []byte
	Network() net.Addr
}

type Converter struct {
	sync.Once
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
		protocol(targetProtocol, targetAddress))
}

func (c *Converter) Convert() (err error) {
	go func() {
		err = c.source.Forward(context.Background(), c.target.Network(), c.source.Packets())
		if err != nil {
			return
		}
	}()

	return
}

func protocol(network string, address string) Protocol {
	switch network {
	case "udp":
		return &Udp{}
	default:
		return &Tcp{}
	}
}
