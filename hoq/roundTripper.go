package hoq

type RounderTripper interface {
	RoundTrip(*Request) (*Response, error)
}

type TCPRounderTripper struct {
}

func (t *TCPRounderTripper) RoundTrip(*Request) (*Response, error) {
	panic("implements me")
	return nil, nil
}

type QUICRounderTripper struct {
}

func (t *QUICRounderTripper) RoundTrip(*Request) (*Response, error) {
	panic("implements me")
	return nil, nil
}
