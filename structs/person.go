package structs

import "fmt"

// ImportBodyMetrics will import extra metrics to a person
func (person *PersonMetrics) ImportBodyMetrics(metrics []BodyMetric) {
	for _, update := range metrics {
		if _, ok := person.BodyMetrics[update.Timestamp]; !ok {
			person.BodyMetrics[update.Timestamp] = BodyMetric{}
		}

		bodyMetric := person.BodyMetrics[update.Timestamp]

		bodyMetric.Weight = update.Weight
		bodyMetric.Timestamp = update.Timestamp
		bodyMetric.Kcal = update.Kcal
		bodyMetric.Fat = update.Fat
		bodyMetric.Tbw = update.Tbw
		bodyMetric.Muscle = update.Muscle
		bodyMetric.Bone = update.Bone
		bodyMetric.Bmi = update.Bmi

		person.BodyMetrics[update.Timestamp] = bodyMetric
	}
}

func (person *PersonMetrics) LastMetric() *BodyMetric {
	var (
		l  BodyMetric
		ts int
	)

	if len(person.BodyMetrics) == 0 {
		fmt.Println("X1")
		return nil
	}

	for i, p := range person.BodyMetrics {
		if i <= ts {
			continue
		}

		ts = i
		l = p
	}

	return &l
}
