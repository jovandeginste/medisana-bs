package structs

// ImportBodyMetrics will import extra metrics to a person
func (person *PersonMetrics) ImportBodyMetrics(metrics []BodyMetric) {
	for _, update := range metrics {
		_, ok := person.BodyMetrics[update.Timestamp]
		if !ok {
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
