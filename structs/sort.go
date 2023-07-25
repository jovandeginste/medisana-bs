package structs

func (s BodyMetrics) Len() int {
	return len(s)
}

func (s BodyMetrics) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s BodyMetrics) Less(i, j int) bool {
	return s[i].Timestamp < s[j].Timestamp
}
