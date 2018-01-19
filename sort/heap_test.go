package sort

import (
	"testing"
	"fmt"
)

func TestHeap(t *testing.T) {
	a := IntSlice{}
	list := []int{}
	for i := 10; i >=0; i -- {
		a = Heappush(a, i)
		list = append(list, i)
	}
	h := IntSlice(list)
	Heapify(h)
	fmt.Println(h, list) // [0 1 4 2 6 5 8 3 7 9 10] [0 1 4 2 6 5 8 3 7 9 10]
	a, val := Heappop(a)
	fmt.Println(a, val) // [1 2 5 4 3 9 6 10 7 8] 0
	Heappushpop(a, 0)
	fmt.Println(a) // [1 2 5 4 3 9 6 10 7 8]
	Heappushpop(a, 11)
	fmt.Println(a) // [2 3 5 4 8 9 6 10 7 11]
}
