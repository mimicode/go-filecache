package go_filecache

import (
	"fmt"
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

func init() {

}
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
	fmt.Println(string(bytes))
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
