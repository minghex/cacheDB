package cache

//统计缓存数据大小
type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (this *Stat) add(k string, v []byte) {
	this.Count += 1
	this.KeySize += int64(len(k))
	this.ValueSize += int64(len(v))
}

func (this *Stat) del(k string, v []byte) {
	this.Count -= 1
	this.KeySize -= int64(len(k))
	this.ValueSize -= int64(len(v))
}
