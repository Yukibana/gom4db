/*
@Time : 2019/12/13 16:46
@Author : Minus4
*/
package rpc

import (
	"fmt"
	"golang.org/x/net/context"
	"gom4db/cache"
	"gom4db/network/cluster"
	"gom4db/network/protobuf"
)

type CacheServiceImpl struct {
	cache cache.KeyValueCache
	cluster.Node
}

func (c *CacheServiceImpl) Get(ctx context.Context, pKey *protobuf.String) (req *protobuf.Response, err error) {
	key := pKey.Value
	addr, ok := c.ShouldProcess(key)
	req = new(protobuf.Response)
	if !ok {
		req.RedirectDir = addr
		req.ErrorMsg = fmt.Sprintf("Redirct: %s", addr)
		return req, nil
	}
	value, err := c.cache.Get(key)

	if value != nil {
		req.Ok = true
		req.Value = cache.Bytes2str(value)
	} else {
		req.ErrorMsg = fmt.Sprintf("The value of key {%s} not found", key)
	}
	return
}

func (c *CacheServiceImpl) Set(ctx context.Context, param *protobuf.SetParam) (req *protobuf.Response, err error) {
	key := param.Key
	addr, ok := c.ShouldProcess(key)
	req = new(protobuf.Response)
	if !ok {
		req.ErrorMsg = fmt.Sprintf("Redirct: %s", addr)
		req.RedirectDir = addr
		return req, nil
	}
	value := param.Value
	err = c.cache.Set(key, value)
	if err != nil {
		req.ErrorMsg = fmt.Sprintf("Err: %s", err)
	} else {
		req.Ok = true
	}
	return
}

func (c *CacheServiceImpl) Del(ctx context.Context, pKey *protobuf.String) (req *protobuf.Response, err error) {
	key := pKey.Value
	addr, ok := c.ShouldProcess(key)
	req = new(protobuf.Response)
	if !ok {
		req.RedirectDir = addr
		req.ErrorMsg = fmt.Sprintf("Redirct: %s", addr)
		return req, nil
	}
	err = c.cache.Del(key)
	if err != nil {
		req.ErrorMsg = fmt.Sprintf("Err: %s", err)
	} else {
		req.Ok = true
	}
	return
}
func NewCacheService(n cluster.Node) protobuf.CacheServiceServer {
	return &CacheServiceImpl{cache.NewKeyValueCache(), n}
}
