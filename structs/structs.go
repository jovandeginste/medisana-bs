package structs

import (
	"time"
)

// PersonMetrics has all data about a single person, including a list of measurements (body metrics)
type PersonMetrics struct {
	Person      int
	Gender      string
	Age         int
	Size        int
	Activity    string
	Updated     bool
	BodyMetrics map[int]BodyMetric
}

// BodyMetrics is shorthand for a list of BodyMetrics
type BodyMetrics []BodyMetric

// BodyMetric is a single tuple of measurements for a given person
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

// AnnotatedBodyMetric contains the values of BodyMetric plus some custom annotations plugins to show
type AnnotatedBodyMetric struct {
	BodyMetric  BodyMetric
	Annotations BodyMetricAnnotations
}

// BodyMetricAnnotations contains annotations to a given BodyMetric
type BodyMetricAnnotations struct {
	Time        time.Time
	DeltaWeight float32
	DeltaFat    float32
	DeltaMuscle float32
	DeltaBone   float32
	DeltaTbw    float32
	DeltaKcal   int
	DeltaBmi    float32
}

// Person contains some fairly static data about a person
type Person struct {
	Valid    bool
	Person   int
	Gender   string
	Age      int
	Size     int
	Activity string
}

// Weight contains a single weight measurement for a person
type Weight struct {
	Valid     bool
	Weight    float32
	Timestamp int
	Person    int
}

// Body contains a secondary measurements for a person
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

// PartialMetric contains either type of measurement sent by the scale
type PartialMetric struct {
	Person Person
	Weight Weight
	Body   Body
}

// Plugin interface describes what a plugin should implement
type Plugin interface {
	Initialize(c Config) Plugin
	ParseData(person *PersonMetrics) bool
}

// Config contains the configuration for the application
type Config struct {
	Device       string
	ScanDuration duration
	DeviceID     string
	Sub          duration
	CsvDir       string
	TimeOffset   int
	Fakeit       bool
	Plugins      map[string]PluginConfig
}

// PluginConfig contains any possible Plugin configuration
type PluginConfig struct {
	Server        string
	SenderName    string
	SenderAddress string
	TemplateFile  string
	Subject       string
	Metrics       int
	Dir           string
	StartTLS      bool
	Recipients    map[string]MailRecipient
}

// MailRecipient contains a person's name and a list of mail addresses to get updates
type MailRecipient struct {
	Name    string
	Address []string
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
