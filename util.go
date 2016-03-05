package consulKvMigrator

import consul "github.com/hashicorp/consul/api"

func compare(source []simpleKv, target consul.KVPairs) differences {
	targetMap := map[string]*consul.KVPair{}
	for idx := range target {
		targetMap[target[idx].Key] = target[idx]
	}

	diffs := differences{}

	// compare values
	for idx := range source {
		sourceKv := source[idx]
		if targetKv, ok := targetMap[sourceKv.key]; ok {
			targetKvValue := string(targetKv.Value)
			if sourceKv.value != targetKvValue {
				diffs = append(diffs, difference{
					key:            sourceKv.key,
					sourceValue:    sourceKv.value,
					targetValue:    targetKvValue,
					targetOriginal: targetKv,
				})
			}
		} else {
			diffs = append(diffs, difference{
				key:         sourceKv.key,
				sourceValue: sourceKv.value,
			})
		}
	}
	return diffs
}

type differences []difference

type difference struct {
	key            string
	sourceValue    string
	targetValue    string
	targetOriginal *consul.KVPair
}
