package plugins

import (
	"bytes"
	"fmt"
	"github.com/jovandeginste/medisana-bs/structs"
	"html/template"
	"log"
	"net/smtp"
	"sort"
	"strconv"
	"strings"
	"time"
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
	Name    string
	Address []string
}

// Initialize the Mail plugin
func (plugin Mail) Initialize(c structs.Config) structs.Plugin {
	newc := c.Plugins["mail"]
	p := Mail{newc.Server, newc.SenderName, newc.SenderAddress, newc.Recipients, newc.StartTLS,
		newc.TemplateFile, newc.Subject, newc.Metrics}
	plugin = p
	log.Println("[PLUGIN MAIL] I am the Mail plugin")
	log.Printf("[PLUGIN MAIL]   - Server: %s\n", plugin.Server)
	log.Printf("[PLUGIN MAIL]   - StartTLS: %t [ To be implemented !!! ]\n", plugin.StartTLS)
	log.Printf("[PLUGIN MAIL]   - SenderName: %s\n", plugin.SenderName)
	log.Printf("[PLUGIN MAIL]   - SenderAddress: %s\n", plugin.SenderAddress)
	log.Printf("[PLUGIN MAIL]   - TemplateFile: %s\n", plugin.TemplateFile)
	log.Printf("[PLUGIN MAIL]   - Subject: %s\n", plugin.Subject)
	log.Printf("[PLUGIN MAIL]   - Metrics: %d\n", plugin.Metrics)
	log.Printf("[PLUGIN MAIL]   - Recipients: %d\n", len(plugin.Recipients))
	return plugin
}

// ParseData will parse new data for a given person
func (plugin Mail) ParseData(person *structs.PersonMetrics) bool {
	log.Println("[PLUGIN MAIL] The mail plugin is parsing new data")
	plugin.sendMail(person)
	return true
}

func (plugin Mail) sendMail(person *structs.PersonMetrics) {
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
	lastMetrics := make(map[time.Time]structs.BodyMetric)
	for _, value := range metrics[len(metrics)-plugin.Metrics:] {
		thetime := time.Unix(int64(value.Timestamp), 0)
		lastMetrics[thetime] = value
	}

	from := fmt.Sprintf("\"%s\" <%s>", plugin.SenderName, plugin.SenderAddress)
	to := recipient.Address

	var auth smtp.Auth

	var msg string
	msg = msg + fmt.Sprintf("From: %s\n", from)
	msg = msg + fmt.Sprintf("To: %s\n", strings.Join(to, ";"))
	msg = msg + fmt.Sprintf("Subject: %s\n", subject)
	msg = msg + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"

	msg = msg + "\n"
	parameters := struct {
		Name     string
		PersonID int
		Metrics  map[time.Time]structs.BodyMetric
	}{
		Name:     recipient.Name,
		PersonID: personID,
		Metrics:  lastMetrics,
	}
	body, err := parseTemplate(plugin.TemplateFile, parameters)
	if err != nil {
		log.Printf("[PLUGIN MAIL] An error occurred: %+v\n", err)
		return
	}
	msg = msg + body

	log.Printf("[PLUGIN MAIL] Sending mail from %s to %s...\n", from, to)
	smtp.SendMail(plugin.Server, auth, plugin.SenderAddress, to, []byte(msg))
	log.Printf("[PLUGIN MAIL] Message was %d bytes.\n", len(msg))
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
