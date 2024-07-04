package cemux

import (
	"context"
	"net/url"

	"anime.bike/hrd"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol"
)

type ResponseWriter interface {
	WriteEvent(*event.Event)
	WriteResult(protocol.Result)
}

type Request struct {
	Ctx   context.Context
	Event event.Event
	Url   *url.URL
}

type Handler interface {
	ServeEvent(*Request, ResponseWriter)
}

type HandlerFunc func(*Request, ResponseWriter)

func (h HandlerFunc) ServeEvent(r *Request, rw ResponseWriter) {
	h(r, rw)
}

type Route struct {
	Pattern string
	Handler Handler
}

type route struct {
	MatcherFunc func(r *Request) bool
	Handler     Handler
}

type Middleware func(Handler) Handler

type Mux struct {
	handler Handler

	mx          []Middleware
	routes      []route
	matcherFunc func(r *Request) bool

	children []*Mux
}

func NewMux() *Mux {
	o := &Mux{}
	o.build()
	return o
}

func (m *Mux) Use(mw ...Middleware) {
	m.mx = append(m.mx, mw...)
	m.build()
}

func (m *Mux) HandleFunc(pattern string, handler HandlerFunc) {
	m.Handle(pattern, handler)
}
func (m *Mux) Handle(pattern string, handler Handler) {
	tokens := hrd.Tokenize([]byte(pattern))
	var matcherFunc func(r *Request) bool
	if pattern != "" {
		matcherFunc = func(r *Request) bool {
			return hrd.Matches(tokens, hrd.Subject(r.Event.Type()))
		}
	}
	r := route{
		Handler:     handler,
		MatcherFunc: matcherFunc,
	}
	m.routes = append(m.routes, r)
	m.build()
}

func (m *Mux) Route(pattern string, fn func(m *Mux)) *Mux {
	sr := NewMux()
	tokens := hrd.Tokenize([]byte(pattern))
	if pattern != "" {
		sr.matcherFunc = func(r *Request) bool {
			return hrd.Matches(tokens, hrd.Subject(r.Event.Type()))
		}
	}
	fn(sr)
	m.children = append(m.children, sr)
	m.build()
	return sr
}

func (mux *Mux) ServeEvent(r *Request, rw ResponseWriter) {
	mux.handler.ServeEvent(r, rw)
}

func (mux *Mux) build() {
	handler := Handler(HandlerFunc(func(r *Request, rw ResponseWriter) {
		if mux.matcherFunc != nil && !mux.matcherFunc(r) {
			return
		}
		for _, v := range mux.routes {
			if v.MatcherFunc != nil && !v.MatcherFunc(r) {
				continue
			}
			v.Handler.ServeEvent(r, rw)
		}
		for _, v := range mux.children {
			v.handler.ServeEvent(r, rw)
		}
	}))
	for i := len(mux.mx) - 1; i >= 0; i-- {
		handler = mux.mx[i](handler)
	}
	mux.handler = handler
}
