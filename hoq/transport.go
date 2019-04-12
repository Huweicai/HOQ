package hoq

/**
底层运输载体，目前支持tcp , quic 两种
*/
type transporter interface {
	Listen(addr string) error
	Accept() (Channel, error)
	Close() error
}

type tcpTransporter struct {
}

func (t *tcpTransporter) Listen(addr string) error {
	panic("implement me")
}

func (t *tcpTransporter) Accept() (Channel, error) {
	panic("implement me")
}

func (t *tcpTransporter) Close() error {
	panic("implement me")
}
