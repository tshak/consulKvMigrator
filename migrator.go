package consulKvMigrator

import consul "github.com/hashicorp/consul/api"

// Migrate migrates kv pairs from a specified migration file to consul
func Migrate(consulKvClient consulKvClient, filename string, prompt bool, dryRun bool) error {
	migrations := parseFromFile(filename)
	fromConsul, err := getFromConsul(consulKvClient, migrations)
	if err != nil {
		return err
	}

	handler := buildHandlers(consulKvClient, prompt, dryRun)

	return handler.handle(compare(migrations, fromConsul))
}

func getFromConsul(consulKvClient consulKvClient, simpleKvs []simpleKv) (consul.KVPairs, error) {
	fromConsul := consul.KVPairs{}
	for idx := range simpleKvs {
		kvPair, _, err := consulKvClient.Get(simpleKvs[idx].key, &consul.QueryOptions{RequireConsistent: true})
		if err != nil {
			return nil, err
		}

		if kvPair != nil {
			fromConsul = append(fromConsul, kvPair)
		}
	}

	return fromConsul, nil
}
