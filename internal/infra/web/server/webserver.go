package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HandlerFuncMethod struct {
	HandleFunc http.HandlerFunc
	Method     string
}

type WebServer struct {
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
	RequestLimit  int
}

func NewWebServer(serverPort string, requestLimit int) *WebServer {
	return &WebServer{
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
		RequestLimit:  requestLimit,
	}
}

func (s *WebServer) AddRoute(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() error {
	mux := http.NewServeMux()
	for path, handler := range s.Handlers {
		mux.Handle(path, s.rateLimitMiddleware(http.HandlerFunc(handler)))
	}

	return http.ListenAndServe(s.WebServerPort, mux)
}

type output struct {
	Message string
}

func (s *WebServer) rateLimitMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		key := r.Header.Get("API_KEY")
		if key == "" {
			output := &output{
				Message: "API_KEY not accept",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(output)
			return
		}
		fmt.Printf("Rate limit verify, max:%d ip:%s, key:%s\n", s.RequestLimit, ip, key)
		handler.ServeHTTP(w, r)
	}
}
