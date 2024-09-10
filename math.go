package splitter

import (
	"github.com/chewxy/math32"
)

func mean(vs [][]float32) []float32 {
	var result []float32

	for i := range vs[0] {
		var sum float32

		for j := range vs {
			sum += vs[j][i]
		}

		result = append(result, sum/float32(len(vs)))
	}

	return result
}

func norm(v []float32) float32 {
	var sum float32

	for _, v := range v {
		sum += v * v
	}

	return math32.Sqrt(sum)
}

func dot(a, b []float32) float32 {
	var sum float32

	for i := range a {
		sum += a[i] * b[i]
	}

	return sum
}
