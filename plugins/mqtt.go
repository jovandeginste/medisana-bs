package plugins

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jovandeginste/medisana-bs/structs"
)

var measurements = []map[string]string{
	{"ha_value": "weight", "scale_value": "weight", "name": "Weight", "icon": "scale-bathroom", "unit": "kg", "class": "weight"},
	{"ha_value": "calories", "scale_value": "kcal", "name": "Calories", "icon": "fire", "unit": "kcal"},
	{"ha_value": "fat", "scale_value": "fat", "name": "Fat", "icon": "account-group", "unit": "%"},
	{"ha_value": "water", "scale_value": "tbw", "name": "Water Ratio", "icon": "water-opacity", "unit": "%"},
	{"ha_value": "muscle", "scale_value": "muscle", "name": "Muscle Ratio", "icon": "weight-lifter", "unit": "%"},
	{"ha_value": "bone", "scale_value": "bone", "name": "Bone Mass", "icon": "bone", "unit": "kg", "class": "weight"},
	{"ha_value": "bmi", "scale_value": "bmi", "name": "BMI", "icon": "calculator-variant-outline", "unit": ""},
}

type MQTT struct {
	Host     string
	Username string
	Password string

	model  string
	client mqtt.Client
}

// Initialize the Csv plugin
func (plugin MQTT) Initialize(c structs.Config) structs.Plugin {
	newc := c.Plugins["mqtt"]
	p := MQTT{
		model:    c.Device,
		Host:     newc.Host,
		Username: newc.Username,
		Password: newc.Password,
	}
	p.initializeClient()

	plugin = p

	log.Println("[PLUGIN MQTT] I am the MQTT plugin")
	log.Printf("[PLUGIN MQTT]   - Model: %s\n", plugin.model)
	log.Printf("[PLUGIN MQTT]   - Host: %s\n", plugin.Host)
	log.Printf("[PLUGIN MQTT]   - Username: %s\n", plugin.Username)

	return plugin
}

func (plugin *MQTT) initializeClient() {
	plugin.client = mqtt.NewClient(
		mqtt.NewClientOptions().
			AddBroker(plugin.Host).
			SetUsername(plugin.Username).
			SetPassword(plugin.Password),
	)
}

type deviceStruct struct {
	Model        string   `json:"mdl"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"mf"`
	Identifiers  []string `json:"identifiers"`
}

type payload struct {
	Name              string       `json:"name"`
	ValueTemplate     string       `json:"value_template"`
	UnitOfMeasurement string       `json:"unit_of_measurement"`
	Icon              string       `json:"icon"`
	StateTopic        string       `json:"state_topic"`
	ObjectID          string       `json:"object_id"`
	UniqueID          string       `json:"unique_id"`
	Device            deviceStruct `json:"device"`
	StateClass        string       `json:"state_class"`
	DeviceClass       string       `json:"device_class,omitempty"`
}

func (plugin MQTT) broadcastAutoDiscover(person *structs.PersonMetrics) error {
	identifier := fmt.Sprintf("%s_person_%d", plugin.model, person.Person)
	identifierLower := strings.ToLower(identifier)

	for _, measurement := range measurements {
		measurementIdentifier := fmt.Sprintf("%s_%s", identifierLower, measurement["ha_value"])
		device := deviceStruct{
			Model:        plugin.model,
			Name:         identifierLower,
			Manufacturer: "Medisana",
			Identifiers:  []string{identifier, identifierLower},
		}

		adTopic := fmt.Sprintf("homeassistant/sensor/%s/%s/config", identifierLower, measurement["ha_value"])

		adPayload := payload{
			Name:              measurement["name"],
			ValueTemplate:     fmt.Sprintf("{{ value_json.%s }}", measurement["scale_value"]),
			UnitOfMeasurement: measurement["unit"],
			Icon:              "mdi:" + measurement["icon"],
			StateTopic:        fmt.Sprintf("homeassistant/sensor/%s/state", identifierLower),
			ObjectID:          measurementIdentifier,
			UniqueID:          measurementIdentifier,
			Device:            device,
			StateClass:        "measurement",
		}

		if val, ok := measurement["class"]; ok {
			adPayload.DeviceClass = val
		}

		j, err := json.Marshal(adPayload)
		if err != nil {
			return err
		}

		log.Printf("[PLUGIN MQTT] Publishing Auto Discovery for %s to %s", measurement["scale_value"], adTopic)

		if token := plugin.client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		defer plugin.client.Disconnect(250)

		if token := plugin.client.Publish(adTopic, 1, false, j); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	return nil
}

func (plugin MQTT) ParseData(person *structs.PersonMetrics) bool {
	if err := plugin.broadcastAutoDiscover(person); err != nil {
		log.Printf("[PLUGIN MQTT] Error: %s", err)
		return false
	}

	if err := plugin.sendLastMetric(person); err != nil {
		log.Printf("[PLUGIN MQTT] Error: %s", err)
		return false
	}

	return true
}

func (plugin MQTT) sendLastMetric(person *structs.PersonMetrics) error {
	identifier := fmt.Sprintf("%s_person_%d", plugin.model, person.Person)
	identifierLower := strings.ToLower(identifier)
	adTopic := fmt.Sprintf("homeassistant/sensor/%s/state", identifierLower)

	lastMetric := person.LastMetric()
	if lastMetric == nil {
		return nil
	}

	j, err := json.Marshal(*lastMetric)
	if err != nil {
		return err
	}

	log.Printf("[PLUGIN MQTT] Publishing measurement for %s to %s", identifier, adTopic)

	if token := plugin.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	defer plugin.client.Disconnect(250)

	if token := plugin.client.Publish(adTopic, 1, false, j); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
