/*
@Time : 2019/11/24 13:18
@Author : Minus4
@Software: GoLand
*/
package gonet

import (
	"container/heap"
	"testing"
)

var h *ResponseHeap

func init() {
	h = new(ResponseHeap)
	heap.Init(h)
}
func TestAll(t *testing.T) {
}
