package go_filecache

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"
)

type FileCache struct {
	mutex   sync.Mutex
	config  FileCacheConfig //配置
	sp      string  //文件夹分隔符
	isGcing bool   //是否正在gc
}

func New(config FileCacheConfig) *FileCache {
	return &FileCache{
		mutex:  sync.Mutex{},
		config: config,
		sp:     string(os.PathSeparator),
	}
}

func Default() *FileCache {
	return &FileCache{
		mutex: sync.Mutex{},
		config: FileCacheConfig{
			PathLevel:       2,
			CachePath:       "runtime/cache",
			CacheFileSuffix: ".bin",
			keyPrefix:       "",
			GcProbability:   10,
			DirMode:         0755,
		},
		sp: string(os.PathSeparator),
	}
}

//缓存文件名
func (fc *FileCache) buildKey(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return fc.config.keyPrefix + hex.EncodeToString(h.Sum(nil))
}

//获取缓存文件路径
func (fc *FileCache) getCacheFile(key string) string {
	if fc.config.PathLevel > 0 {
		prexLen := len(fc.config.keyPrefix)
		base := fc.config.CachePath
		for i := 0; i < fc.config.PathLevel; i++ {
			if i*2 <= len(key)-2-prexLen {
				base = base + fc.sp + key[prexLen+i*2:i*2+2+prexLen]
			}
		}
		return path.Join(base, key+fc.config.CacheFileSuffix)
	} else {
		return path.Join(fc.config.CachePath, key+fc.config.CacheFileSuffix)
	}
}

//判断key是否存在
func (fc *FileCache) Exists(key string) bool {
	cacheFile := fc.getCacheFile(fc.buildKey(key))
	fileInfo, err := os.Stat(cacheFile)
	if os.IsNotExist(err) {
		return false
	}
	//最后修改时间 在当前时间之前 说明已过期
	if fileInfo.ModTime().Before(time.Now()) {
		return false
	}
	return true
}

//删除缓存
func (fc *FileCache) Del(key string) {
	fc.mutex.Lock()
	cacheFile := fc.getCacheFile(fc.buildKey(key))
	_ = os.Remove(cacheFile)
	fc.mutex.Unlock()
}

//设置缓存
func (fc *FileCache) Set(key, value string, exp int64) bool {
	cacheFile := fc.getCacheFile(fc.buildKey(key))
	//if fc.config.PathLevel > 0 {
	//	if err := os.MkdirAll(path.Dir(cacheFile), fc.config.DirMode); err != nil {
	//		return false
	//	}
	//}

	if err := os.MkdirAll(path.Dir(cacheFile), fc.config.DirMode); err != nil {
		return false
	}
	fc.mutex.Lock()
	err := ioutil.WriteFile(cacheFile, []byte(value), fc.config.DirMode)
	if exp <= 0 {
		exp = 24 * 3600 * 365 * 2 //2年
	}
	_ = os.Chtimes(cacheFile, time.Now(), time.Now().Add(time.Duration(exp)*time.Second))
	fc.mutex.Unlock()
	if err != nil {
		return false
	}
	return true
}

func (fc *FileCache) Get(key string) []byte {
	cacheFile := fc.getCacheFile(fc.buildKey(key))
	fc.mutex.Lock()
	val, err := ioutil.ReadFile(cacheFile)
	fc.mutex.Unlock()
	if err != nil {
		return nil
	} else {
		return val
	}
}

func (fc *FileCache) gc() {

	n := rand.Intn(1000000)
	//百万分只10的改路发生gc
	if n < fc.config.GcProbability {

		fc.mutex.Lock()
		if fc.isGcing {
			fc.mutex.Unlock()
			return
		}
		fc.isGcing = true
		fc.mutex.Unlock()
		//执行gc
		go fc.gcRecursive(fc.config.CachePath)
	}
}
//递归gc过期文件
func (fc *FileCache) gcRecursive(cachePath string) {
	if dirs, err := ioutil.ReadDir(cachePath); err != nil {
		if cachePath == fc.config.CachePath {
			fc.mutex.Lock()
			fc.isGcing = false
			fc.mutex.Unlock()
		}
	} else {
		for _, dir := range dirs {
			if dir.IsDir() {
				fc.gcRecursive(path.Join(cachePath, dir.Name()))
			} else {
				//fmt.Println(path.Join(cachePath, dir.Name()))
				if dir.ModTime().Before(time.Now()) {
					_ = os.Remove(path.Join(cachePath, dir.Name()))
				}
			}
		}
		if cachePath == fc.config.CachePath {
			fc.mutex.Lock()
			fc.isGcing = false
			fc.mutex.Unlock()
		}
	}

}

//配置
type FileCacheConfig struct {
	//缓存文件目录级别
	PathLevel int
	//缓存目录
	CachePath string
	//缓存文件后缀
	CacheFileSuffix string
	//缓存Key前缀
	keyPrefix string
	//gc 概率 百万分只1  0 - 1000000
	GcProbability int
	//目录权限
	DirMode os.FileMode
}
