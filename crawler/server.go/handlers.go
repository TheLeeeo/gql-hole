package crawlserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TheLeeeo/gql-test-suite/crawler"
	"github.com/TheLeeeo/gql-test-suite/introspection"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) Crawl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ops, err := s.crawler.Crawl()
	if err != nil {
		if err == introspection.ErrNoTargetAddr {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "no target graphql endpoint specified")
			return
		}

		log.Println("error crawling: ", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error crawling: ", err)
		return
	}

	failedOps := make([]crawler.CrawlOperation, 0)
	for _, op := range ops {
		if op.Error != nil || !op.Denied {
			failedOps = append(failedOps, op)
		}
	}

	b, err := json.Marshal(failedOps)
	if err != nil {
		log.Println("error marshalling results of operations: ", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error marshalling results of operations")
	}

	w.Write(b)
}

func (s *Server) GetIgnore(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := json.Marshal(s.crawler.GetIgnore())
	if err != nil {
		log.Println("error marshalling ignore list: ", err)

		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}

func (s *Server) SetIgnore(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ignore []string
	err := json.NewDecoder(r.Body).Decode(&ignore)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "error decoding ignore list, please format as an array ([\"value1\", \"value2\"]))")
		return
	}

	s.crawler.SetIgnore(ignore)
	log.Println("Updated ignore list to: ", ignore)

	fmt.Fprint(w, ignore)
}

func (s *Server) SetTargetURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newUrl string
	if err := json.NewDecoder(r.Body).Decode(&newUrl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "error decoding target URL, please format as a string (\"value\")")
		return
	}

	if err := s.crawler.SetTargetURL(newUrl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "error setting target URL: ", err)
		return
	}

	log.Println("Updated target URL to: ", newUrl)

	fmt.Fprint(w, newUrl)
}

func (s *Server) GetTargetURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, s.crawler.GetTargetURL())
}
