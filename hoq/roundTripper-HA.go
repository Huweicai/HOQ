package hoq

import "github.com/lucas-clemente/quic-go/h2quic"

type QUICCourierHA struct {
	rt *h2quic.RoundTripper
}

func (q *QUICCourierHA) RoundTrip(req *Request) (resp *Response, remote *remoteInfo, err error) {
	r, err := req.wrap()
	if err != nil {
		return
	}
	rr, err := q.rt.RoundTrip(r)
	if err != nil {
		return
	}
	resp, err = NewResponse(rr.StatusCode, convertHttpHeader(rr.Header), rr.Body)
	if err != nil {
		return
	}
	remote = &remoteInfo{}
	return
}
