package main

type PersonMetrics struct {
	Person      int
	Gender      string
	Age         int
	Size        int
	Activity    string
	BodyMetrics map[int]BodyMetrics
}

type BodyMetric struct {
	Timestamp int
	Weight    float32
	Fat       float32
	Muscle    float32
	Bone      float32
	Tbw       float32
	Kcal      int
	Bmi       float32
}
type BodyMetrics []BodyMetric

type Person struct {
	Valid    bool
	Person   int
	Gender   string
	Age      int
	Size     int
	Activity string
}

type Weight struct {
	Valid     bool
	Weight    float32
	Timestamp int
	Person    int
}

type Body struct {
	Valid     bool
	Timestamp int
	Person    int
	Kcal      int
	Fat       float32
	Tbw       float32
	Muscle    float32
	Bone      float32
}
type PartialMetric struct {
	Person Person
	Weight Weight
	Body   Body
}
