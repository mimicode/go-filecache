package main

import (
	"fmt"
	gofilecache "github.com/mimicode/go-filecache"
)

func main() {
	cache := gofilecache.Default()
	cache.Set("xiaoming", "xiaoming", 300)
	val := cache.Get("xiaoming")
	fmt.Println(string(val))
}
