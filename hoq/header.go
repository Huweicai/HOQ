package hoq

import "net/textproto"

/**
头部字段，不应该仅仅只是一个map
针对Host等字段需要单独定义
*/
type Headers struct {
	headers map[string]string
}

func (h *Headers) Serialize() string {
	out := ""
	//todo check omit values after the first one
	for name, value := range h.headers {
		if value == "" {
			continue
		}
		out += name + ": " + value
	}
	return out
}

func (h *Headers) Set(k, v string) {

}

func (h *Headers) Get(k string) string {
	//todo implements it
	return ""
}

func mimeHeaderToMap(mime textproto.MIMEHeader) map[string]string {
	if len(mime) == 0 {
		return nil
	}
	m := make(map[string]string)
	for key, values := range mime {
		if len(values) == 0 {
			continue
		}
		m[key] = values[0]
	}
	return m
}
