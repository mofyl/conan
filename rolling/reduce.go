package rolling

func Sum(iterator Iterator) float64 {
	res := 0.0

	var b Bucket
	for iterator.Next() {
		b = iterator.Bucket()
		for _, v := range b.Point {
			res += v
		}
	}

	return res
}

func Avg(iterator Iterator) float64 {
	res := 0.0
	count := 0.0
	var b Bucket
	for iterator.Next() {
		b = iterator.Bucket()
		for _, v := range b.Point {
			res += v
			count++
		}
	}

	return res / count
}

func Min(iterator Iterator) float64 {
	res := 0.0
	var b Bucket
	for iterator.Next() {
		b = iterator.Bucket()
		for _, v := range b.Point {
			if v < res {
				res = v
			}
		}
	}

	return res
}

func Max(iterator Iterator) float64 {
	res := 0.0
	var b Bucket

	for iterator.Next() {
		b = iterator.Bucket()
		for _, v := range b.Point {
			if v > res {
				res = v
			}
		}
	}

	return res
}

func Count(iterator Iterator) float64 {
	var res int64
	var b Bucket
	for iterator.Next() {
		b = iterator.Bucket()
		res += b.Count
	}
	return float64(res)
}
