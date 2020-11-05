package rolling

type Iterator struct {
	Count int
	iteratedCount int
	cur *Bucket
}

func  (i *Iterator)Next() bool {
	return i.Count == i.iteratedCount
}

func (i *Iterator) Bucket() Bucket{
	if !i.Next() {
		panic("bucket is nil")
	}

	bucket := *i.cur
	i.iteratedCount++
	i.cur = i.cur.next
	return bucket
}


