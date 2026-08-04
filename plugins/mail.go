package plugins

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/smtp"
	"sort"
	"strconv"
	"strings"

	"github.com/jovandeginste/medisana-bs/structs"
)

// Mail contains configuration for the Mail plugin
type Mail struct {
	Server        string
	SenderName    string
	SenderAddress string
	Recipients    map[string]structs.MailRecipient
	StartTLS      bool
	TemplateFile  string
	Subject       string
	Metrics       int
}

// MailRecipient contains a person name, and a list of addresses to send updates to
type MailRecipient struct {
	Address []string
}

func (plugin Mail) Name() string {
	return "MAIL"
}

func (plugin Mail) Logger() *slog.Logger {
	return slog.With("plugin", plugin.Name())
}

// Initialize the Mail plugin
func (plugin Mail) Initialize(c structs.Config) structs.Plugin {
	newc := c.Plugins["mail"]
	plugin = Mail{
		newc.Server, newc.SenderName, newc.SenderAddress, newc.Recipients, newc.StartTLS,
		newc.TemplateFile, newc.Subject, newc.Metrics,
	}

	plugin.Logger().Debug("I am the Mail plugin")
	plugin.Logger().Debug("Config", "server", plugin.Server, "start_tls", plugin.StartTLS, "sender_name", plugin.SenderName, "sender_address", plugin.SenderAddress, "template_file", plugin.TemplateFile, "subject", plugin.Subject, "metrics", plugin.Metrics, "recipients", len(plugin.Recipients))

	return plugin
}

// ParseData will parse new data for a given person
func (plugin Mail) ParseData(person *structs.PersonMetrics) bool {
	plugin.Logger().Info("The mail plugin is parsing new data")

	plugin.sendMail(person)

	return true
}

func (plugin Mail) sendMail(person *structs.PersonMetrics) { //nolint:funlen
	personID := person.Person
	recipient := plugin.Recipients[strconv.Itoa(personID)]
	subject := plugin.Subject

	metrics := make(structs.BodyMetrics, len(person.BodyMetrics))

	idx := 0

	for _, value := range person.BodyMetrics {
		metrics[idx] = value
		idx++
	}

	sort.Sort(metrics)

	lastMetrics := make(map[int]structs.AnnotatedBodyMetric)

	l := len(metrics) - plugin.Metrics - 1
	if l < 0 {
		l = 0
	}

	previousMetric := metrics[l]

	for _, value := range metrics[l+1:] {
		annotations := structs.BodyMetricAnnotations{
			Time:        value.ToTime(),
			DeltaWeight: value.Weight - previousMetric.Weight,
			DeltaFat:    value.Fat - previousMetric.Fat,
			DeltaMuscle: value.Muscle - previousMetric.Muscle,
			DeltaBone:   value.Bone - previousMetric.Bone,
			DeltaTbw:    value.Tbw - previousMetric.Tbw,
			DeltaKcal:   value.Kcal - previousMetric.Kcal,
			DeltaBmi:    value.Bmi - previousMetric.Bmi,
		}
		lastMetrics[value.Timestamp] = structs.AnnotatedBodyMetric{
			Annotations: annotations,
			BodyMetric:  value,
		}

		previousMetric = value
	}

	from := fmt.Sprintf("\"%s\" <%s>", plugin.SenderName, plugin.SenderAddress)
	to := recipient.Address

	var auth smtp.Auth

	msg := strings.Builder{}

	msg.WriteString("From: " + from + "\n")
	msg.WriteString("To: " + strings.Join(to, ",") + "\n")
	msg.WriteString("Subject: " + subject + "\n")
	msg.WriteString("MIME-version: 1.0;\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n")
	msg.WriteString("\n")

	parameters := struct {
		PersonName string
		PersonID   int
		Metrics    map[int]structs.AnnotatedBodyMetric
	}{
		PersonName: person.Name,
		PersonID:   personID,
		Metrics:    lastMetrics,
	}

	body, err := parseTemplate(plugin.TemplateFile, parameters)
	if err != nil {
		plugin.Logger().Error("Template parsing failed", "err", err)
		return
	}

	msg.WriteString(body)

	plugin.Logger().Info("Sending mail", "from", from, "to", to)

	if err := smtp.SendMail(plugin.Server, auth, plugin.SenderAddress, to, []byte(msg.String())); err != nil {
		plugin.Logger().Error("Sending mail failed", "err", err)
		return
	}

	plugin.Logger().Debug("Message sent", "bytes", msg.Len())
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	var result string

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return result, err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return result, err
	}

	result = buf.String()

	return result, nil
}

func (plugin Mail) InitializeData(_ *structs.PersonMetrics) bool {
	return true
}
