package tcpdump

import "context"

type CustomDecoder struct {
	name       string
	decodeFunc func(ctx context.Context, r Reader, p *TcpPacket) (Message, error)
}

func (c *CustomDecoder) Decode(ctx context.Context, r Reader, p *TcpPacket) (Message, error) {
	return c.decodeFunc(ctx, r, p)
}

func (c *CustomDecoder) Name() string {
	return c.name
}

var _ Decoder = (*CustomDecoder)(nil)

func NewCustomDecoder(name string, decodeFunc func(ctx context.Context, r Reader, p *TcpPacket) (Message, error)) *CustomDecoder {
	if name == "" || decodeFunc == nil {
		panic("invalid custom decoder params")
	}
	return &CustomDecoder{
		name:       name,
		decodeFunc: decodeFunc,
	}
}
