package propagator

import "net/http"

type Propagator struct {
	cookieName string
	cookieOpt  func(*http.Cookie)
}

func NewPropagator() *Propagator {
	return &Propagator{
		cookieName: "cookie",
		cookieOpt: func(cookie *http.Cookie) {

		},
	}
}

func (p *Propagator) Inject(writer http.ResponseWriter, id string) error {
	cookie := &http.Cookie{Name: p.cookieName, Value: id}
	p.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter, id string) error {
	cookie := &http.Cookie{Name: p.cookieName, Value: id, MaxAge: -1}
	http.SetCookie(writer, cookie)
	return nil
}
