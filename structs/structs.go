package structs

import (
	"time"
)

type PersonMetrics struct {
	Person      int
	Gender      string
	Age         int
	Size        int
	Activity    string
	Updated     bool
	BodyMetrics map[int]BodyMetric
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

type Plugins map[string]Plugin

type Plugin interface {
	Initialize() bool
	ParseData(person *PersonMetrics) bool
}

type Config struct {
	Device       string
	ScanDuration duration
	DeviceID     string
	Sub          duration
	CsvDir       string
	Time_offset  int
	Fakeit       bool
	Plugins      interface{}
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
func (d duration) AsTimeDuration() time.Duration {
	return d.Duration
}
