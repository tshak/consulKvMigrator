package consulKvMigrator

import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/hashicorp/consul/api"

func TestCompare_Equal(t *testing.T) {
	fromFile := []simpleKv{
		simpleKv{
			key:   "key1",
			value: "value1",
		},
	}

	fromConsul := api.KVPairs{
		&api.KVPair{
			Key:   "key1",
			Value: []byte("value1"),
		},
	}

	result := compare(fromFile, fromConsul)
	assert.Equal(t, 0, len(result))
}

func TestCompare_NotEqual(t *testing.T) {
	fromFile := []simpleKv{
		createSimpleKvChanged("1"), // changed
		createSimpleKvChanged("2"), // added
	}

	fromConsul := api.KVPairs{
		createConsulKv("1"),
		createConsulKv("3"), // ignore (we don't delete uknown targets)
	}

	result := compare(fromFile, fromConsul)
	assert.Equal(t, 2, len(result))

	diff1 := result[0]
	assert.Equal(t, "key1", diff1.key)
	assert.Equal(t, "value1_changed", diff1.sourceValue)
	assert.Equal(t, "value1", diff1.targetValue)
	assert.Equal(t, fromConsul[0], diff1.targetOriginal)

	diff2 := result[1]
	assert.Equal(t, "key2", diff2.key)
	assert.Equal(t, "value2_changed", diff2.sourceValue)
}

func createSimpleKvChanged(id string) simpleKv {
	return simpleKv{key: "key" + id, value: "value" + id + "_changed"}
}

func createConsulKv(id string) *api.KVPair {
	return &api.KVPair{
		Key:   "key" + id,
		Value: []byte("value" + id),
	}
}
