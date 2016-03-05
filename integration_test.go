// +build integration

package consulKvMigrator

import "fmt"
import "flag"
import "testing"
import "github.com/stretchr/testify/assert"
import consul "github.com/hashicorp/consul/api"

const baseKeyPath = "consulKvMigratorIntegrationTest"

var address = flag.String("address", "localhost:8500", "address for consul")

func TestSimpleMigration(t *testing.T) {
	consulConfig := &consul.Config{
		Address: *address,
	}
	client, err := consul.NewClient(consulConfig)
	if err != nil {
		t.Error(err)
	}

	helper := consulTestHelper{t: t, client: client}

	defer helper.cleanupConsul()

	existingKvPair := consul.KVPair{Key: helper.createKey("existingKey"), Value: []byte("original")}

	if _, err := client.KV().Put(&existingKvPair, nil); err != nil {
		t.Error(err)
	}

	if err := Migrate(client.KV(), "integration_test.json", false, false); err != nil {
		t.Error(err)
	}

	helper.assertKVValueEquals("existingKey", "existingKeyValue")
	helper.assertKVValueEquals("newKey", "newKeyValue")
}

type consulTestHelper struct {
	t      *testing.T
	client *consul.Client
}

func (h *consulTestHelper) createKey(key string) string {
	return fmt.Sprintf("%s/%s", baseKeyPath, key)
}

func (h *consulTestHelper) cleanupConsul() {
	h.client.KV().DeleteTree(baseKeyPath, nil)
}

func (h *consulTestHelper) assertKVValueEquals(expectedKey string, expectedValue string) {
	fullKey := h.createKey(expectedKey)
	kv, _, err := h.client.KV().Get(fullKey, nil)
	if err != nil {
		h.t.Error(err)
	}

	if kv == nil {
		h.t.Errorf("Missing key in consul: '%s'", fullKey)
	}

	assert.Equal(h.t, expectedValue, string(kv.Value))
}
