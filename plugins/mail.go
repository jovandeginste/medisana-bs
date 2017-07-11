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

type Mail struct {
	Server        string
	SenderName    string
	SenderAddress string
	Recipients    map[int]MailRecipient
	StartTLS      bool
	TemplateFile  string
	Subject       string
	Metrics       int
}
type MailRecipient struct {
	Name    string
	Address []string
}

func (plugin Mail) Initialize(c interface{}) structs.Plugin {
	newc := c.(map[string]interface{})
	plugin.Server = newc["Server"].(string)
	plugin.SenderName = newc["SenderName"].(string)
	plugin.SenderAddress = newc["SenderAddress"].(string)
	if newc["StartTLS"] != nil {
		plugin.StartTLS = newc["StartTLS"].(bool)
	} else {
		plugin.StartTLS = false
	}
	plugin.TemplateFile = newc["TemplateFile"].(string)
	plugin.Subject = newc["Subject"].(string)
	plugin.Metrics = int(newc["Metrics"].(int64))
	plugin.Recipients = make(map[int]MailRecipient)

	recipients := newc["Recipients"].(map[string]interface{})
	for number, rec_config := range recipients {
		id, err := strconv.Atoi(number)
		if err != nil {
			log.Printf("[PLUGIN MAIL] Configuration for non-numerical recipient id '%v'", id)
		} else {
			conv_rec_config := rec_config.(map[string]interface{})
			unconv_address := conv_rec_config["Address"].([]interface{})
			address := make([]string, len(unconv_address))
			for i := range unconv_address {
				address[i] = unconv_address[i].(string)
			}
			log.Printf("%+v", address)
			plugin.Recipients[id] = MailRecipient{
				Name:    conv_rec_config["Name"].(string),
				Address: address,
			}
		}
	}

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
func (plugin Mail) ParseData(person *structs.PersonMetrics) bool {
	log.Println("[PLUGIN MAIL] The mail plugin is parsing new data")
	plugin.sendMail(person)
	return true
}
func (mail Mail) sendMail(person *structs.PersonMetrics) {
	personId := person.Person
	recipient := mail.Recipients[personId]
	subject := mail.Subject

	metrics := make(structs.BodyMetrics, len(person.BodyMetrics))
	idx := 0
	for _, value := range person.BodyMetrics {
		metrics[idx] = value
		idx++
	}
	sort.Sort(metrics)
	lastMetrics := make(map[time.Time]structs.BodyMetric)
	for _, value := range metrics[len(metrics)-mail.Metrics:] {
		thetime := time.Unix(int64(value.Timestamp), 0)
		lastMetrics[thetime] = value
	}

	from := fmt.Sprintf("\"%s\" <%s>", mail.SenderName, mail.SenderAddress)
	to := recipient.Address

	var auth smtp.Auth
	auth = nil

	var msg string
	msg = msg + fmt.Sprintf("From: %s\n", from)
	msg = msg + fmt.Sprintf("To: %s\n", strings.Join(to, ";"))
	msg = msg + fmt.Sprintf("Subject: %s\n", subject)
	msg = msg + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"

	msg = msg + "\n"
	parameters := struct {
		Name     string
		PersonId int
		Metrics  map[time.Time]structs.BodyMetric
	}{
		Name:     recipient.Name,
		PersonId: personId,
		Metrics:  lastMetrics,
	}
	body, err := ParseTemplate(mail.TemplateFile, parameters)
	if err != nil {
		log.Printf("[PLUGIN MAIL] An error occurred: %+v\n", err)
		return
	}
	msg = msg + body

	log.Printf("[PLUGIN MAIL] Sending mail from %s to %s...\n", from, to)
	smtp.SendMail(mail.Server, auth, mail.SenderAddress, to, []byte(msg))
	log.Printf("[PLUGIN MAIL] Message was %d bytes.\n", len(msg))
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
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
