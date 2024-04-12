package payloadcms_test

import "github.com/ainsleydev/webkit/pkg/apis/payloadcms"

func Newwww() {
	p := payloadcms.New("http://localhost", "1234")

	p.Get()
}
