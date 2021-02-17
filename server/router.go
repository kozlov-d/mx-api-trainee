package server

import (
	"log"
	"net/http"
	"time"
)

//Route describes each route
type Route struct {
	Method  string
	Path    string
	Queries []string
	Handler http.HandlerFunc
}

func contentTypeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")
		h.ServeHTTP(w, r)
	})
}

func loggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func (s *Server) setupRouter() {
	s.Router.Use(loggerMiddleware)
	s.Router.Use(contentTypeMiddleware)
	var routes = []Route{
		{
			"PUT",
			"/merchants/{id:[0-9]+}",
			[]string{},
			s.createTask,
		},
		{
			"GET",
			"/tasks/{id:[0-9]+}",
			[]string{},
			s.getTaskByID,
		},
		{
			"GET",
			"/offers",
			[]string{"offerId", "{offerID:[0-9]+}", "merchantId", "{merchantID:[0-9]+}", "sub", "{sub:[A-Za-z0-9_\\s]+}"},
			s.getOffers,
		},
		{
			"GET",
			"/offers",
			[]string{"offerId", "{offerID:[0-9]+}", "merchantId", "{merchantID:[0-9]+}"},
			s.getOffers,
		},
		{
			"GET",
			"/offers",
			[]string{"merchantId", "{merchantID:[0-9]+}", "sub", "{sub:[A-Za-z0-9_\\s]+}"},
			s.getOffers,
		},
		{
			"GET",
			"/offers",
			[]string{"offerId", "{offerID:[0-9]+}", "sub", "{sub:[A-Za-z0-9_\\s]+}"},
			s.getOffers,
		},
		{
			"GET",
			"/offers",
			[]string{"offerId", "{offerID:[0-9]+}"},
			s.getOffers,
		},
		{
			"GET",
			"/offers",
			[]string{"merchantId", "{merchantID:[0-9]+}"},
			s.getOffers,
		},
		{
			"GET",
			"/offers",
			[]string{"sub", "{sub:[A-Za-z0-9_\\s]+}"},
			s.getOffers,
		},
		{
			"GET",
			"/offers",
			[]string{},
			s.getOffers,
		},
	}
	for _, route := range routes {
		var handler http.Handler
		handler = route.Handler
		s.Router.
			Methods(route.Method).
			PathPrefix(route.Path).
			Queries(route.Queries...).
			Handler(handler)
	}
}
