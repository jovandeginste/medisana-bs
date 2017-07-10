package structs

import (
	"sort"
)

func (s BodyMetrics) Len() int {
	return len(s)
}
func (s BodyMetrics) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s BodyMetrics) Less(i, j int) bool {
	return s[i].Timestamp < s[j].Timestamp
}
func mergeSort(w1 BodyMetrics, w2 BodyMetrics) (result BodyMetrics) {
	var exists bool
	for _, v1 := range w1 {
		result = append(result, v1)
	}

	for _, v2 := range w2 {
		exists = false
		for _, v1 := range w1 {
			if v1.Timestamp == v2.Timestamp {
				exists = true
			}
		}
		if !exists {
			result = append(result, v2)
		}
	}
	sort.Sort(result)
	return
}
