package plugins

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jovandeginste/medisana-bs/structs"
)

var measurements = []map[string]string{
	{"ha_value": "name", "scale_value": "name", "name": "Name", "icon": "card-account-details-outline", "unit": "", "state_class": ""},
	{"ha_value": "weight", "scale_value": "weight", "name": "Weight", "icon": "scale-bathroom", "unit": "kg", "class": "weight"},
	{"ha_value": "timestamp", "scale_value": "timestamp", "name": "Last measurement", "icon": "calendar-clock-outline", "unit": "", "class": "timestamp"},
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

func (plugin MQTT) Name() string {
	return "MQTT"
}

func (plugin MQTT) Logger() *slog.Logger {
	return slog.With("plugin", plugin.Name())
}

// Initialize the MQTT plugin
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

	plugin.Logger().Debug("I am the MQTT plugin")
	plugin.Logger().Debug("Config", "model", plugin.model, "host", plugin.Host, "username", plugin.Username)

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

func (plugin MQTT) ensureConnected() error {
	if plugin.client.IsConnected() && plugin.client.IsConnectionOpen() {
		return nil
	}

	if token := plugin.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
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
	UnitOfMeasurement string       `json:"unit_of_measurement,omitempty"`
	Icon              string       `json:"icon"`
	StateTopic        string       `json:"state_topic"`
	ObjectID          string       `json:"object_id"`
	UniqueID          string       `json:"unique_id"`
	Device            deviceStruct `json:"device"`
	StateClass        string       `json:"state_class,omitempty"`
	DeviceClass       string       `json:"device_class,omitempty"`
}

func (plugin MQTT) broadcastAutoDiscover(person *structs.PersonMetrics) error {
	identifier := strings.ToLower(fmt.Sprintf("%s_person_%d", plugin.model, person.Person))

	if err := plugin.ensureConnected(); err != nil {
		return err
	}

	for _, measurement := range measurements {
		measurementIdentifier := fmt.Sprintf("%s_%s", identifier, measurement["ha_value"])
		device := deviceStruct{
			Model:        plugin.model,
			Name:         identifier,
			Manufacturer: "Medisana",
			Identifiers:  []string{identifier},
		}

		adTopic := fmt.Sprintf("homeassistant/sensor/%s/%s/config", identifier, measurement["ha_value"])

		adPayload := payload{
			Name:              measurement["name"],
			ValueTemplate:     fmt.Sprintf("{{ value_json.%s }}", measurement["scale_value"]),
			UnitOfMeasurement: measurement["unit"],
			Icon:              "mdi:" + measurement["icon"],
			StateTopic:        fmt.Sprintf("homeassistant/sensor/%s/state", identifier),
			ObjectID:          measurementIdentifier,
			UniqueID:          measurementIdentifier,
			Device:            device,
		}

		if val, ok := measurement["state_class"]; ok {
			adPayload.StateClass = val
		} else {
			adPayload.StateClass = "measurement"
		}

		if val, ok := measurement["class"]; ok {
			adPayload.DeviceClass = val
		}

		j, err := json.Marshal(adPayload)
		if err != nil {
			return err
		}

		plugin.Logger().Debug("Publishing auto discovery", "measurement", measurement["scale_value"], "topic", adTopic)
		plugin.Logger().Debug("Payload", "payload", string(j))

		if token := plugin.client.Publish(adTopic, 1, false, j); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	return nil
}

func (plugin MQTT) ParseData(person *structs.PersonMetrics) bool {
	plugin.Logger().Info("The MQTT plugin is parsing new data")

	if err := plugin.sendLastMetric(person); err != nil {
		plugin.Logger().Error("Error sending last metric", "err", err)
		return false
	}

	return true
}

func (plugin MQTT) sendLastMetric(person *structs.PersonMetrics) error {
	identifier := strings.ToLower(fmt.Sprintf("%s_person_%d", plugin.model, person.Person))
	adTopic := fmt.Sprintf("homeassistant/sensor/%s/state", identifier)

	lastMetric := person.LastMetric()
	if lastMetric == nil {
		return nil
	}

	j, err := json.Marshal(*lastMetric)
	if err != nil {
		return err
	}

	plugin.Logger().Info("Publishing measurement", "identifier", identifier, "topic", adTopic)
	plugin.Logger().Info("Payload", "payload", string(j))

	if err := plugin.ensureConnected(); err != nil {
		return err
	}

	if token := plugin.client.Publish(adTopic, 1, false, j); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (plugin MQTT) InitializeData(person *structs.PersonMetrics) bool {
	plugin.Logger().Info("The MQTT plugin is initializing the last data", "person_id", person.Person, "name", person.Name)

	if err := plugin.broadcastAutoDiscover(person); err != nil {
		plugin.Logger().Error("Error broadcasting auto discovery", "err", err)
		return false
	}

	go func() {
		for {
			if err := plugin.sendLastMetric(person); err != nil {
				plugin.Logger().Error("Error sending periodic metric", "err", err)
			}

			time.Sleep(300 * time.Second)
		}
	}()

	return true
}
