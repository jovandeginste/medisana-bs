package main

import (
	"log"
	"math"
	"time"
)

var allPersons = make([]*PersonMetrics, 8)

func MetricParser() {
	for i := range allPersons {
		allPersons[i] = &PersonMetrics{Person: i + 1, BodyMetrics: make(map[int]BodyMetric)}
		allPersons[i].ImportBodyMetrics(ImportCsv(i + 1))
	}
	sync_chan := make(chan bool)
	Debounce(3*time.Second, sync_chan)
	for {
		partial_metric := <-metric_chan
		UpdatePerson(partial_metric.Person)
		UpdateBody(partial_metric.Body)
		UpdateWeight(partial_metric.Weight)
		sync_chan <- true
	}
}

func GetPersonMetrics(personId int) *PersonMetrics {
	person := allPersons[personId-1]
	return person
}

func UpdatePerson(update Person) {
	if !update.Valid {
		return
	}
	log.Printf("Received person metrics: %+v", update)
	person := GetPersonMetrics(update.Person)
	person.Gender = update.Gender
	person.Age = update.Age
	person.Size = update.Size
	person.Activity = update.Activity
	PrintPerson(person)
}

func UpdateBody(update Body) {
	if !update.Valid {
		return
	}
	log.Printf("Received body metrics: %+v", update)
	person := GetPersonMetrics(update.Person)
	person.Updated = true
	_, ok := person.BodyMetrics[update.Timestamp]
	if !ok {
		log.Printf("No body metric - creating")
		person.BodyMetrics[update.Timestamp] = BodyMetric{}
	}
	bodyMetric := person.BodyMetrics[update.Timestamp]
	bodyMetric.Timestamp = update.Timestamp
	bodyMetric.Kcal = update.Kcal
	bodyMetric.Fat = update.Fat
	bodyMetric.Tbw = update.Tbw
	bodyMetric.Muscle = update.Muscle
	bodyMetric.Bone = update.Bone
	person.BodyMetrics[update.Timestamp] = bodyMetric
	PrintPerson(person)
}
func UpdateWeight(update Weight) {
	if !update.Valid {
		return
	}
	log.Printf("Received weight metrics: %+v", update)
	person := GetPersonMetrics(update.Person)
	person.Updated = true
	_, ok := person.BodyMetrics[update.Timestamp]
	if !ok {
		log.Printf("No body metric - creating")
		person.BodyMetrics[update.Timestamp] = BodyMetric{}
	}
	bodyMetric := person.BodyMetrics[update.Timestamp]
	bodyMetric.Weight = update.Weight
	bodyMetric.Timestamp = update.Timestamp
	if bodyMetric.Weight > 0 && person.Size > 0 {
		bodyMetric.Bmi = bodyMetric.Weight / float32(math.Pow(float64(person.Size)/100, 2))
	}

	person.BodyMetrics[update.Timestamp] = bodyMetric
	PrintPerson(person)
}
func PrintPerson(person *PersonMetrics) {
	log.Printf("Person %d now has %d metrics.\n", person.Person, len(person.BodyMetrics))
}

func Debounce(lull time.Duration, in chan bool) {
	go func() {
		for {
			select {
			case <-in:
			case <-time.Tick(lull):
				for _, person := range allPersons {
					if person.Updated {
						log.Printf("Person %d was updated.\n", person.Person)
						person.ExportBodyMetrics()
					}
				}
			}
		}
	}()
}
