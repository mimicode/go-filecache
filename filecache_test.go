package go_filecache

import (
	"fmt"
	"strconv"
	"testing"
)

var (
	filecache = Default()
	fnew      = New(FileCacheConfig{
		PathLevel:       10,
		CachePath:       "runtime/cache",
		CacheFileSuffix: ".bbbb",
		keyPrefix:       "cache_",
		GcProbability:   100,
	})
)

func TestFileCache_buildKey(t *testing.T) {
	key := filecache.buildKey("afda")
	fmt.Println(key)
	key = fnew.buildKey("afda")
	fmt.Println(key)
}

func TestFileCache_getCacheFile(t *testing.T) {
	file := filecache.getCacheFile(filecache.buildKey("afda"))
	fmt.Println(file)
	file = fnew.getCacheFile(fnew.buildKey("afda"))
	fmt.Println(file)
}

func TestFileCache_Set(t *testing.T) {
	filecache.Set("12212", "dada", 30)
}

func TestFileCache_Get(t *testing.T) {
	bytes := filecache.Get("12212")
	fmt.Println(bytes)
}

func TestFileCache_Exists(t *testing.T) {
	exists := filecache.Exists("12212")
	fmt.Println(exists)
}

func TestFileCache_Del(t *testing.T) {
	filecache.Del("12212")
}

func TestFileCache_gcRecursive(t *testing.T) {
	filecache.gcRecursive(filecache.config.CachePath)
}

func BenchmarkFileCache_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filecache.Set("V"+"_"+strconv.Itoa(i), "1", 1000)
	}
}
func BenchmarkFileCache_Set20Tread(b *testing.B) {
	q := make(chan bool, 0)
	cn := 20
	for c := 0; c < cn; c++ {
		go func() {
			for i := 0; i < b.N; i++ {
				filecache.Set("V"+strconv.Itoa(cn)+"_"+strconv.Itoa(i), "1", 1000)
			}
			q <- true
		}()
	}
	for c := 0; c < cn; c++ {
		<-q
	}
}
func BenchmarkFileCache_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filecache.Get("V" + "_" + strconv.Itoa(i))
	}
}
func BenchmarkFileCache_Get20Thread(b *testing.B) {
	q := make(chan bool, 0)
	cn := 20
	for c := 0; c < cn; c++ {
		go func() {
			for i := 0; i < b.N; i++ {
				filecache.Get("V" + strconv.Itoa(cn) + "_" + strconv.Itoa(i))
			}
			q <- true
		}()
	}
	for c := 0; c < cn; c++ {
		<-q
	}
}

func BenchmarkFileCache_SetGet(b *testing.B) {
	ch := make(chan struct{}, 0)
	go func() {
		for i := 0; i < b.N; i++ {
			filecache.Set("V"+"_"+strconv.Itoa(i), "1", 1000)
		}
		ch <- struct{}{}
	}()

	go func() {
		for i := 0; i < b.N; i++ {
			filecache.Get("V" + "_" + strconv.Itoa(i))
		}
		ch <- struct{}{}
	}()

	<-ch
	<-ch

}
