package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/minghex/cacheDB/benchmark/cacheClient"
)

//统计时间消耗
type statistic struct {
	count int           //桶内操作的数量
	time  time.Duration //桶内操作总耗时
}

type result struct {
	getCount    int
	setCount    int
	missCount   int
	statBuckets []statistic
}

func (r *result) addStatistic(bucket int, stat statistic) {
	if bucket >= len(r.statBuckets) {
		nStatBuckets := make([]statistic, bucket+1)
		copy(nStatBuckets, r.statBuckets)
		r.statBuckets = nStatBuckets
	}

	r.statBuckets[bucket].count += stat.count
	r.statBuckets[bucket].time += stat.time
}

func (r *result) addDuration(d time.Duration, resultType string) {
	//根据消耗的时间分桶
	b := int(d / time.Millisecond)
	r.addStatistic(b, statistic{1, d})

	if resultType == "get" {
		r.getCount++
	} else if resultType == "set" {
		r.setCount++
	} else {
		r.missCount++
	}
}

func (r *result) addResult(src *result) {
	for b, s := range src.statBuckets {
		r.addStatistic(b, s)
	}

	r.getCount += src.getCount
	r.setCount += src.setCount
	r.missCount += src.missCount
}

func run(client cacheClient.Client, cmd *cacheClient.Cmd, r *result) {
	//get 操作期望值
	expect := cmd.Value
	//记录操作消耗时间
	start := time.Now()
	client.Run(cmd)
	d := time.Now().Sub(start)
	//操作类型
	resultType := cmd.OpName
	if resultType == "get" {
		if cmd.Value == "" {
			resultType = "miss"
		} else if cmd.Value != expect {
			panic(cmd)
		}
	}

	r.addDuration(d, resultType)
}

func pipeline(client cacheClient.Client, cmds []*cacheClient.Cmd, r *result) {
	//校验Get操作
	expects := make([]string, len(cmds))
	for k, c := range cmds {
		if c.OpName == "get" {
			expects[k] = c.Value
		}
	}

	start := time.Now()
	client.PipelineRun(cmds)
	d := time.Now().Sub(start)

	for k, c := range cmds {
		resultType := c.OpName
		if c.OpName == "get" {
			if c.Value == "" {
				resultType = "miss"
			} else if c.Value != expects[k] {
				panic(c.Value)
			}
		}
		r.addDuration(d, resultType)
	}
}

//初始化全局变量
var port, threads, total, valueSize, keyspacelen, pipelen int
var typ, addr, operation string

func init() {
	flag.StringVar(&typ, "type", "redis", "cache type")
	flag.StringVar(&addr, "h", "localhost", "cache server address")
	flag.IntVar(&port, "p", 12345, "cache server port")
	flag.StringVar(&operation, "t", "set", "test set ,could be set/get/mixed")
	flag.IntVar(&threads, "c", 1, "number of parallel connections")
	flag.IntVar(&total, "n", 1000, "number of requests")
	flag.IntVar(&valueSize, "d", 1000, "data size of SET/GET value in bytes")
	flag.IntVar(&pipelen, "P", 1, "pipeline length")
	flag.IntVar(&keyspacelen, "r", 12345, "keyspacelen, use random keys from 0 to keyspacelen-1")
	flag.Parse()
	fmt.Println("type is", typ)
	fmt.Println("server addr is", addr)
	fmt.Println("server port is", port)
	fmt.Println("total", total, "requests")
	fmt.Println("data size is", valueSize)
	fmt.Println("we have", threads, "connections")
	fmt.Println("operation is", operation)
	fmt.Println("keyspacelen is", keyspacelen)
	fmt.Println("pipeline length is", pipelen)
}

func operate(id, count int, ch chan *result) {
	client := cacheClient.New(typ, addr)
	cmds := make([]*cacheClient.Cmd, 0)
	//根据valuesize 创建一个value值
	valuePrefix := strings.Repeat("a", valueSize)
	r := &result{0, 0, 0, make([]statistic, 0)}
	for i := 0; i < count; i++ {
		var tmp int
		if keyspacelen > 0 {
			tmp = rand.Intn(keyspacelen)
		} else {
			tmp = id*count + i
		}
		key := fmt.Sprintf("%d", tmp)
		value := fmt.Sprintf("%s%d", valuePrefix, tmp)
		name := operation
		if operation == "mixed" {
			if rand.Intn(2) == 1 {
				name = "set"
			} else {
				name = "get"
			}
		}

		c := &cacheClient.Cmd{
			OpName: name,
			Key:    key,
			Value:  value,
			Error:  nil,
		}

		if pipelen > 1 {
			cmds = append(cmds, c)
			if len(cmds) == pipelen {
				pipeline(client, cmds, r)
				cmds = make([]*cacheClient.Cmd, 0)
			}
		} else {
			run(client, c, r)
		}
	}

	if len(cmds) != 0 {
		pipeline(client, cmds, r)
	}

	ch <- r
}

func main() {
	ch := make(chan *result, threads)
	res := &result{0, 0, 0, make([]statistic, 0)}
	start := time.Now()
	for i := 0; i < threads; i++ {
		go operate(i, total/threads, ch)
	}
	for i := 0; i < threads; i++ {
		res.addResult(<-ch)
	}

	//输出benchmark信息
	d := time.Now().Sub(start)
	totalCount := res.getCount + res.missCount + res.setCount
	fmt.Printf("%d records get\n", res.getCount)
	fmt.Printf("%d records miss\n", res.missCount)
	fmt.Printf("%d records set\n", res.setCount)
	fmt.Printf("%f seconds total\n", d.Seconds())

	statCountSum := 0
	statTimeSum := time.Duration(0)
	for b, s := range res.statBuckets {
		if s.count == 0 {
			continue
		}
		statTimeSum += s.time
		statCountSum += s.count
		fmt.Printf("%d%% requests <= %d ms \n", statCountSum*100/totalCount, b)
	}

	//每个请求平均消耗时长
	fmt.Printf("%d usec average for each request\n", int64(statTimeSum/time.Microsecond)/int64(statCountSum))
	//单位时间吞吐量
	fmt.Printf("throughput is %f MB/s\n", float64((res.getCount+res.setCount)*valueSize)/1e6/d.Seconds())
	//单位时间响应数
	fmt.Printf("rps is %f\n", float64(totalCount)/float64(d.Seconds()))
}
