package router

import "HOQ/hoq"

func defaultHeader() *hoq.Headers {
	return hoq.NewHeaders(map[string]string{"Server": ServerName})
}
