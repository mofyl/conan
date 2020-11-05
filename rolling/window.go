package rolling

type Bucket struct {
	Point []float64
	Count int64
	next  *Bucket
}

func (b *Bucket) Append(val float64) {
	b.Point = append(b.Point, val)
	b.Count++
}

func (b *Bucket) Add(offset int, val float64) {
	b.Point[offset] += val
	b.Count++
}

func (b *Bucket) Reset() {
	b.Point = b.Point[:0]
	b.Count = 0
}

func (b *Bucket) Next() *Bucket {
	return b.next
}

type Window struct {
	window []Bucket
	size   int
}

type WindowOpt struct {
	Size int
}

func NewWindow(opt WindowOpt) *Window {
	bucket := make([]Bucket, opt.Size)

	nextOffset := 0
	for i := 0; i < opt.Size; i++ {

		bucket[i] = Bucket{
			Point: make([]float64, 0),
		}
		nextOffset = i + 1

		if nextOffset == opt.Size {
			nextOffset = 0
		}
		bucket[i].next = &bucket[nextOffset]
	}

	return &Window{window: bucket, size: opt.Size}
}

func (w *Window) ResetWindow() {
	for offset := range w.window {
		w.ResetBucket(offset)
	}
}
func (w *Window) ResetBucket(offset int) {
	w.window[offset].Reset()
}

func (w *Window) ResetBucketFromOffset(offset []int) {
	for _, v := range offset {
		w.ResetBucket(v)
	}
}

func (w *Window) Append(offset int, val float64) {
	w.window[offset].Append(val)
}

// 这里现在默认是向 Point 0 号位置去Add
func (w *Window) Add(offset int, val float64) {

	if w.window[offset].Count == 0 {
		w.window[offset].Append(val)
		return
	}
	w.window[offset].Add(0, val)
}

func (w *Window) Bucket(offset int) Bucket {

	return w.window[offset]
}

func (w *Window) Size() int {
	return w.size
}

func (w *Window) Iterator(offset int, count int) Iterator {
	return Iterator{
		Count: count,
		cur:   &w.window[offset],
	}
}
