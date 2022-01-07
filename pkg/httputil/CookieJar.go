package httputil

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"net/url"
)

type CookieJar struct {
	cookies map[*url.URL][]*http.Cookie
}

func (j CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.cookies[u] = append(j.cookies[u], cookies...)
}

func (j CookieJar) Cookies(u *url.URL) []*http.Cookie {
	if cs, ok := j.cookies[u]; !ok {
		return []*http.Cookie{}
	} else {
		return cs
	}
}

func (j CookieJar) Encode() []byte {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(j)
	return buf.Bytes()
}

func (j CookieJar) IsZero() bool {
	return len(j.cookies) == 0
}

func Decode(bs []byte) CookieJar {
	dec := gob.NewDecoder(bytes.NewBuffer(bs))
	j := CookieJar{}
	dec.Decode(&j)
	return j
}
