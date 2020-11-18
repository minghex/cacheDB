package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type cacheHandler struct {
	*Server
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}

func (this *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//解析Path
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		value, code, _ := this.get(key)
		w.WriteHeader(code)
		w.Write(value)
	case http.MethodPut:
		code, _ := this.set(key, r)
		w.WriteHeader(code)
	case http.MethodDelete:
		code, _ := this.del(key)
		w.WriteHeader(code)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (this *cacheHandler) get(key string) ([]byte, int, error) {
	v, e := this.Get(key)
	if e != nil {
		return nil, http.StatusNotFound, e
	}
	return v, http.StatusOK, nil
}

func (this *cacheHandler) set(key string, r *http.Request) (int, error) {
	v, _ := ioutil.ReadAll(r.Body)
	if len(v) != 0 {
		e := this.Set(key, v)
		if e != nil {
			return http.StatusInternalServerError, nil
		}
		return http.StatusOK, e
	}
	return http.StatusBadRequest, fmt.Errorf("value is empty")
}

func (this *cacheHandler) del(key string) (int, error) {
	e := this.Del(key)
	if e != nil {
		return http.StatusBadRequest, e
	}
	return http.StatusOK, nil
}
