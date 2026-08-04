package main

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/jovandeginste/medisana-bs/plugins"
	"github.com/jovandeginste/medisana-bs/structs"
)

var allPersons = make([]*structs.PersonMetrics, 8)

func metricParserLogger() *slog.Logger {
	return slog.With("component", "metricparser")
}

// StartMetricParser will initialize the Persons from csv and parse incoming metrics
func StartMetricParser() {
	for name, c := range config.People {
		p := &structs.PersonMetrics{
			Person:      c.ID,
			Name:        name,
			BodyMetrics: make(map[int]structs.BodyMetric),
		}
		p.ImportBodyMetrics(structs.ImportCsv(c.ID))

		metricParserLogger().Debug("Imported person", "person_id", c.ID, "name", p.Name, "metrics", len(p.BodyMetrics))

		plugins.InitializeData(p)

		allPersons[c.ID] = p
	}

	go parseMetrics()
}

func parseMetrics() {
	for {
		partialMetric := <-metricChan

		updatePerson(partialMetric.Person)
		updateBody(partialMetric.Body)
		updateWeight(partialMetric.Weight)
	}
}

func getPersonMetrics(personID int) *structs.PersonMetrics {
	return allPersons[personID]
}

func updatePerson(update structs.Person) {
	if !update.Valid {
		return
	}

	metricParserLogger().Info("Received person metrics", "person", fmt.Sprintf("%+v", update))

	person := getPersonMetrics(update.Person)
	person.Gender = update.Gender
	person.Age = update.Age
	person.Size = update.Size
	person.Activity = update.Activity

	printPerson(person)
}

func updateBody(update structs.Body) {
	if !update.Valid {
		return
	}

	metricParserLogger().Info("Received body metrics", "body", fmt.Sprintf("%+v", update))

	person := getPersonMetrics(update.Person)
	person.Updated = true

	if _, ok := person.BodyMetrics[update.Timestamp]; !ok {
		metricParserLogger().Info("No body metric - creating")

		person.BodyMetrics[update.Timestamp] = structs.BodyMetric{}
	}

	bodyMetric := person.BodyMetrics[update.Timestamp]
	bodyMetric.Timestamp = update.Timestamp
	bodyMetric.Kcal = update.Kcal
	bodyMetric.Fat = update.Fat
	bodyMetric.Tbw = update.Tbw
	bodyMetric.Muscle = update.Muscle
	bodyMetric.Bone = update.Bone
	bodyMetric.TimestampString = update.ToRFC3339()
	bodyMetric.Name = person.Name

	person.BodyMetrics[update.Timestamp] = bodyMetric

	printPerson(person)
}

func updateWeight(update structs.Weight) {
	if !update.Valid {
		return
	}

	metricParserLogger().Info("Received weight metrics", "weight", fmt.Sprintf("%+v", update))

	person := getPersonMetrics(update.Person)
	person.Updated = true

	if _, ok := person.BodyMetrics[update.Timestamp]; !ok {
		metricParserLogger().Info("No body metric - creating")

		person.BodyMetrics[update.Timestamp] = structs.BodyMetric{}
	}

	bodyMetric := person.BodyMetrics[update.Timestamp]
	bodyMetric.Weight = update.Weight
	bodyMetric.Timestamp = update.Timestamp
	bodyMetric.TimestampString = update.ToRFC3339()
	bodyMetric.Name = person.Name

	if bodyMetric.Weight > 0 && person.Size > 0 {
		bodyMetric.Bmi = bodyMetric.Weight / float32(math.Pow(float64(person.Size)/100, 2))
	}

	person.BodyMetrics[update.Timestamp] = bodyMetric

	printPerson(person)
}

func printPerson(person *structs.PersonMetrics) {
	metricParserLogger().Info("Person metrics updated", "person_id", person.Person, "name", person.Name, "metrics", len(person.BodyMetrics))
}

func debounce() {
	for _, person := range allPersons {
		if person == nil || !person.Updated {
			continue
		}

		metricParserLogger().Info("Person updated, calling all plugins", "person_id", person.Person, "name", person.Name)
		plugins.ParseData(person)
		person.Updated = false
	}
}
