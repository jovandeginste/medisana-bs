package main

type BodyMetric struct {
	Timestamp  int
	BodyMetric float32
	Fat        float32
	Muscle     float32
	Bone       float32
	Tbw        float32
	Kcal       int
	Bmi        float32
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
	Valid      bool
	BodyMetric float32
	Timestamp  int
	Person     int
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
