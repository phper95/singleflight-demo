package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/singleflight"
)

var n = int32(0)

var singleFlightGetArticle singleflight.Group

func cache(id int) int32 {
	//fmt.Println("get by cache:", id)
	return atomic.LoadInt32(&n)
}

func db(id int) int32 {
	fmt.Println("get by db:", id)
	return 168
}

func get(id int, key string) int32 {
	_n := cache(id)
	if _n == 0 {
		val, _, _ := singleFlightGetArticle.Do(key, func() (interface{}, error) {
			return db(id), nil
		})
		_n = val.(int32)
		atomic.StoreInt32(&n, _n)
	}
	return _n
}

func main() {
	wg := sync.WaitGroup{}
	key := "demo"

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			get(i, key)
		}()
	}

	wg.Wait()
	fmt.Println("---------------------------")
}

