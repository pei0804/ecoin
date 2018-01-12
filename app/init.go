package app

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pei0804/topicoin/controller"
	"github.com/pei0804/topicoin/view"
)

// Server Server
type Server struct {
	router *chi.Mux
}

// New Server構造体のコンストラクタ
func New() *Server {
	return &Server{
		router: chi.NewRouter(),
	}
}

func (s *Server) Config() {
	// In debug mode, we compile templates on every request.
	view.Init(template.FuncMap{}, true)
}

// Middleware ミドルウェア
func (s *Server) Middleware() {
}

// Router ルーティング設定
func (s *Server) Router() {
	cron := controller.NewCronController()
	s.router.Route("/api", func(api chi.Router) {
		api.Route("/cron", func(c chi.Router) {
			c.Get("/", apiHandler(cron.Show).ServeHTTP)
		})
	})
	s.router.NotFound(notFoundHandler)
}

func init() {
	s := New()
	s.Config()
	s.Middleware()
	s.Router()
	http.Handle("/", s.router)
}
