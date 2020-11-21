package cacheClient

type Cmd struct {
	OpName string //操作类型
	Key    string //操作的键
	Value  string //操作的值
	Error  error
}

type Client interface {
	Run(*Cmd)
	PipelineRun(cmds []*Cmd)
}

func New(typ, addr string) Client {
	if typ == "http" {
		return newHttpClient(addr)
	}

	if typ == "redis" {
		return newRedisClient(addr)
	}

	if typ == "tcp" {
		return newTcpClient(addr)
	}

	panic(typ)
}
