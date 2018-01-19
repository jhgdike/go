package sort

type Array interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func siftdown(heap Array, startPos, pos int) {
	for pos > startPos {
		parentPos := (pos-1)>>1
		if !heap.Less(pos, parentPos) {
			break
		}
		heap.Swap(pos, parentPos)
		pos = parentPos
	}
}

func siftup(heap Array, pos int) {
	startpos := pos
	childpos := pos*2 + 1
	for childpos < heap.Len() {
		rightpos := childpos + 1
		if rightpos < heap.Len() && heap.Less(rightpos, childpos) {
			childpos = rightpos
		}
		heap.Swap(pos, childpos)
		pos = childpos
		childpos = pos*2 + 1
	}
	siftdown(heap, startpos, pos)
}

type IntSlice []int

func (p IntSlice) Len() int {return len(p)}
func (p IntSlice) Less(i, j int) bool {return p[i] < p[j]}
func (p IntSlice) Swap(i, j int) {p[i], p[j] = p[j], p[i]}

func Heappush(heap IntSlice, item int) IntSlice{
	heap = IntSlice(append([]int(heap), item))
	siftdown(heap, 0, heap.Len() - 1)
	return heap
}

func Heappop(heap IntSlice) (IntSlice, int) {
	lastItem, heap := heap[heap.Len()-1], IntSlice([]int(heap)[:heap.Len()-1])
	if heap.Len() > 0 {
		return heap, Heappoppush(heap, lastItem)
	}
	return heap, lastItem
}

func Heappoppush(heap IntSlice, item int) int{
	item, heap[0] = heap[0], item
	siftup(heap, 0)
	return item
}

func Heappushpop(heap IntSlice, item int) int{
	if heap.Len() > 0 && heap[0] < item {
		Heappoppush(heap, item)
	}
	return item
}

func Heapify(heap IntSlice) {
	for i := (heap.Len() / 2) -1; i >= 0; i -- {
		siftup(heap, i)
	}
}
