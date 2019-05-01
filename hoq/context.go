package hoq

import "time"

type Context struct {
	Remote   *remoteInfo
	Request  *Request
	Response *Response
}

func (Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (Context) Done() <-chan struct{} {
	return nil
}

func (Context) Err() error {
	return nil
}

func (Context) Value(key interface{}) interface{} {
	return ""
}
