package main

import (
	"fmt"
	"gom4db/cache"
	"unsafe"
)

func main(){
	c := cache.NewCache()

	_ = c.Set("2", []byte("hello"))
	re,_ := c.Get("2")
	fmt.Println("The result is "+Bytes2str(re))
}
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
