package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/study-go-rate-limit/internal/usecase"
)

type HandlerFuncMethod struct {
	HandleFunc http.HandlerFunc
	Method     string
}

type WebServer struct {
	Handlers         map[string]http.HandlerFunc
	WebServerPort    string
	RateLimitUseCase *usecase.RateLimitUseCase
}

func NewWebServer(serverPort string, rateLimitUseCase *usecase.RateLimitUseCase) *WebServer {
	return &WebServer{
		Handlers:         make(map[string]http.HandlerFunc),
		WebServerPort:    serverPort,
		RateLimitUseCase: rateLimitUseCase,
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

type ErrOut struct {
	Message string
}

func (s *WebServer) rateLimitMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		key := r.Header.Get("API_KEY")

		if ip != "" && key != "" {
			output, err := s.RateLimitUseCase.Execute(r.Context(), usecase.RateLimitUseCaseInputDto{
				Ip:  ip,
				Key: key,
			})

			if err == nil && !output.Allow {
				errMsg := &ErrOut{
					Message: "you have reached the maximum number of requests or actions allowed within a certain time frame",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(errMsg)
				return
			}
		}

		handler.ServeHTTP(w, r)
	}
}
