package gonet

type response struct {
	body  []byte
	order int
}

// An ResponseHeap is a min-heap of integers.
type ResponseHeap []response

func (h ResponseHeap) Len() int {
	return len(h)
}

func (h ResponseHeap) Less(i, j int) bool {
	return h[i].order < h[j].order
}
func (h ResponseHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *ResponseHeap) Push(x interface{}) { *h = append(*h, x.(response)) }

func (h ResponseHeap) IsEmpty() bool {
	return h.Len() == 0
}
func (h *ResponseHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h ResponseHeap) Top() (order int) {

	x := h[0].order
	return x
}
