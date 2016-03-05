package consulKvMigrator

import "fmt"
import "time"
import consul "github.com/hashicorp/consul/api"

type consulHandler struct {
	chainableDifferenceHandler
	consulKvClient consulKvClient
}

func (handler consulHandler) handle(differences differences) error {
	startTime := time.Now()
	for idx := range differences {
		difference := differences[idx]
		if difference.targetOriginal == nil {
			newKVPair := &consul.KVPair{Key: difference.key, Value: []byte(difference.sourceValue)}
			if _, err := handler.consulKvClient.Put(newKVPair, nil); err != nil {
				return err
			}
		} else {
			if err := handler.checkAndSetWithRetry(difference.targetOriginal, []byte(difference.sourceValue)); err != nil {
				return err
			}
		}
	}

	fmt.Printf("Changes migrated: %d\nTook: %.2f seconds\n", len(differences), time.Since(startTime).Seconds())
	return handler.next(differences)
}

func (handler consulHandler) checkAndSetWithRetry(kvPair *consul.KVPair, newValue []byte) error {
	kvPair.Value = newValue
	success, _, err := handler.consulKvClient.CAS(kvPair, nil)
	if err != nil {
		return err
	}

	if !success {
		retryKVPair, _, err := handler.consulKvClient.Get(kvPair.Key, &consul.QueryOptions{RequireConsistent: true})
		if err != nil {
			return err
		}
		return handler.checkAndSetWithRetry(retryKVPair, newValue)
	}

	return nil
}
