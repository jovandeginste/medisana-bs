package main

import (
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jovandeginste/medisana-bs/plugins"
	"github.com/jovandeginste/medisana-bs/structs"
)

var allPersons = make([]*structs.PersonMetrics, 8)

// MetricParser will initialize the Persons from csv and parse incoming metrics
func MetricParser() {
	for name, c := range config.People {
		p := &structs.PersonMetrics{
			Person:      c.ID,
			Name:        name,
			BodyMetrics: make(map[int]structs.BodyMetric),
		}
		p.ImportBodyMetrics(structs.ImportCsv(c.ID))

		log.Debugf("[METRIC PARSER] Imported person %d (%s) with %d metrics", c.ID, p.Name, len(p.BodyMetrics))

		plugins.InitializeData(p)

		allPersons[c.ID] = p
	}

	syncChan := make(chan bool)

	go debounce(3*time.Second, syncChan)

	for {
		partialMetric := <-metricChan

		updatePerson(partialMetric.Person)
		updateBody(partialMetric.Body)
		updateWeight(partialMetric.Weight)

		syncChan <- true
	}
}

func getPersonMetrics(personID int) *structs.PersonMetrics {
	return allPersons[personID]
}

func updatePerson(update structs.Person) {
	if !update.Valid {
		return
	}

	log.Infof("[METRIC PARSER] Received person metrics: %+v", update)

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

	log.Infof("[METRIC PARSER] Received body metrics: %+v", update)

	person := getPersonMetrics(update.Person)
	person.Updated = true

	if _, ok := person.BodyMetrics[update.Timestamp]; !ok {
		log.Infof("[METRIC PARSER] No body metric - creating")

		person.BodyMetrics[update.Timestamp] = structs.BodyMetric{}
	}

	bodyMetric := person.BodyMetrics[update.Timestamp]
	bodyMetric.Timestamp = update.Timestamp
	bodyMetric.Kcal = update.Kcal
	bodyMetric.Fat = update.Fat
	bodyMetric.Tbw = update.Tbw
	bodyMetric.Muscle = update.Muscle
	bodyMetric.Bone = update.Bone
	person.BodyMetrics[update.Timestamp] = bodyMetric

	printPerson(person)
}

func updateWeight(update structs.Weight) {
	if !update.Valid {
		return
	}

	log.Infof("[METRIC PARSER] Received weight metrics: %+v", update)

	person := getPersonMetrics(update.Person)
	person.Updated = true

	if _, ok := person.BodyMetrics[update.Timestamp]; !ok {
		log.Infof("[METRIC PARSER] No body metric - creating")

		person.BodyMetrics[update.Timestamp] = structs.BodyMetric{}
	}

	bodyMetric := person.BodyMetrics[update.Timestamp]
	bodyMetric.Weight = update.Weight
	bodyMetric.Timestamp = update.Timestamp

	if bodyMetric.Weight > 0 && person.Size > 0 {
		bodyMetric.Bmi = bodyMetric.Weight / float32(math.Pow(float64(person.Size)/100, 2))
	}

	person.BodyMetrics[update.Timestamp] = bodyMetric

	printPerson(person)
}

func printPerson(person *structs.PersonMetrics) {
	log.Infof("[METRIC PARSER] Person %d (%s) now has %d metrics.", person.Person, person.Name, len(person.BodyMetrics))
}

func debounce(lull time.Duration, in chan bool) {
	for {
		select {
		case <-in:
		case <-time.Tick(lull):
			for _, person := range allPersons {
				if person == nil || !person.Updated {
					continue
				}

				log.Infof("[METRIC PARSER] Person %d (%s) was updated -- calling all plugins.", person.Person, person.Name)
				plugins.ParseData(person)
				person.Updated = false
			}
		}
	}
}
