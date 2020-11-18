package http

import (
	"encoding/json"
	"net/http"
)

type statusHandler struct {
	*Server
}

func (s *Server) statusHandler() *statusHandler {
	return &statusHandler{s}
}

func (this *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, e := json.Marshal(this.GetStat())
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
