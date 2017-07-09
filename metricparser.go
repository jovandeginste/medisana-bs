package main

import (
	"log"
)

var allPersons = make([]*PersonMetrics, 8)

func MetricParser() {
	for i := range allPersons {
		allPersons[i] = &PersonMetrics{Person: i + 1, BodyMetrics: make(map[int]BodyMetrics)}
	}
	for {
		partial_metric := <-metric_chan
		log.Printf("Received partial metric: %+v\n", partial_metric)
		UpdatePerson(partial_metric.Person)
		UpdateBody(partial_metric.Body)
		UpdateWeight(partial_metric.Weight)
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
	PrintPerson(person)
}
func UpdateWeight(update Weight) {
	if !update.Valid {
		return
	}
	log.Printf("Received weight metrics: %+v", update)
	person := GetPersonMetrics(update.Person)
	bodyMetric := person.BodyMetrics[update.Timestamp]
	log.Printf("Body metric: %+v\n", bodyMetric)
	PrintPerson(person)
}
func PrintPerson(person *PersonMetrics) {
	log.Printf("Person information: %+v", person)
}
