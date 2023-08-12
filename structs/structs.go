package structs

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// PersonMetrics has all data about a single person, including a list of measurements (body metrics)
type PersonMetrics struct {
	Name        string
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
	Timestamp       int     `json:"-"`
	TimestampString string  `json:"timestamp"`
	Weight          float32 `json:"weight"`
	Fat             float32 `json:"fat"`
	Muscle          float32 `json:"muscle"`
	Bone            float32 `json:"bone"`
	Tbw             float32 `json:"tbw"`
	Kcal            int     `json:"kcal"`
	Bmi             float32 `json:"bmi"`
}

func (b *BodyMetric) ToTime() time.Time {
	return time.Unix(int64(b.Timestamp), 0)
}

func (b *BodyMetric) ToRFC3339() string {
	return b.ToTime().Format(time.RFC3339)
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

func (w *Weight) ToTime() time.Time {
	return time.Unix(int64(w.Timestamp), 0)
}

func (w *Weight) ToRFC3339() string {
	return w.ToTime().Format(time.RFC3339)
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

func (b *Body) ToTime() time.Time {
	return time.Unix(int64(b.Timestamp), 0)
}

func (b *Body) ToRFC3339() string {
	return b.ToTime().Format(time.RFC3339)
}

// PartialMetric contains either type of measurement sent by the scale
type PartialMetric struct {
	Person Person
	Weight Weight
	Body   Body
}

// Plugin interface describes what a plugin should implement
type Plugin interface {
	Name() string
	Initialize(c Config) Plugin
	ParseData(person *PersonMetrics) bool
	InitializeData(person *PersonMetrics) bool
	Logger() log.FieldLogger
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
	People       map[string]PersonConfig
	Plugins      map[string]PluginConfig
}

type PersonConfig struct {
	ID int
}

// PluginConfig contains any possible Plugin configuration
type PluginConfig struct {
	Server        string
	SenderName    string
	SenderAddress string
	TemplateFile  string
	Subject       string
	Metrics       int
	StartTLS      bool
	Recipients    map[string]MailRecipient
	Dir           string
	Host          string
	Username      string
	Password      string
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
