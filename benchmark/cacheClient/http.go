package cacheClient

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	addr string
}

func newHttpClient(addr string) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, "http://" + addr + ":12345/cache/"}
}

func (this *httpClient) Run(cmd *Cmd) {
	if cmd.OpName == "get" {
		cmd.Value = this.get(cmd.Key)
		return
	}

	if cmd.OpName == "set" {
		this.set(cmd.Key, cmd.Value)
		return
	}

	panic("unknown cmd name " + cmd.OpName)
}

func (this *httpClient) PipelineRun(cmds []*Cmd) {
	for _, c := range cmds {
		this.Run(c)
	}
}

func (this *httpClient) set(k, v string) {
	req, err := http.NewRequest(http.MethodPut, this.addr+k, strings.NewReader(v))
	if err != nil {
		panic(err)
	}

	resp, err := this.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}

func (this *httpClient) get(key string) string {
	resp, err := this.Get(this.addr + key)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(b)
}
