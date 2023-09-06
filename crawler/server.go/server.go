package crawlserver

import (
	"log"
	"net/http"

	"github.com/TheLeeeo/gql-test-suite/crawler"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	crawler *crawler.Crawler

	cfg Config
}

func New(cfg Config) *Server {
	cr := crawler.New(cfg.CrawlerConfig)

	return &Server{
		crawler: cr,
		cfg:     cfg,
	}
}

func (s *Server) Run() error {
	router := s.SetupRouter()

	s.crawler.StartPolling()

	log.Println("Starting crawl server on ", s.cfg.HttpPort)
	return http.ListenAndServe(s.cfg.HttpPort, router)
}

func (s *Server) SetupRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/crawl", s.Crawl)

	router.GET("/ignore", s.GetIgnore)
	router.POST("/ignore", s.SetIgnore)

	router.GET("/target", s.GetTargetURL)
	router.POST("/target", s.SetTargetURL)

	router.PanicHandler = s.PanicHandler

	return router
}

func (s *Server) PanicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println("Panic: ", err)
}
