package consulKvMigrator

import "testing"
import "github.com/stretchr/testify/assert"
import consul "github.com/hashicorp/consul/api"

func TestHandle_PutsNewItem(t *testing.T) {
	diffs := append(differences{}, difference{key: "myKey", sourceValue: "myValue"})

	consulKvClient := new(mockConsulKvClient)
	handler := consulHandler{consulKvClient: consulKvClient}

	handler.handle(diffs)

	assert.Equal(t, 1, len(consulKvClient.Puts))
	assert.Equal(t, "myKey", consulKvClient.Puts[0].Key)
	assert.Equal(t, []byte("myValue"), consulKvClient.Puts[0].Value)
}

func TestHandle_SetExistingItem(t *testing.T) {
	diffs := append(differences{}, difference{
		key:            "myKey",
		sourceValue:    "myValue",
		targetValue:    "oldValue",
		targetOriginal: &consul.KVPair{Flags: 2, Key: "myKey", Value: []byte("oldValue"), CreateIndex: 2, ModifyIndex: 3},
	})

	consulKvClient := &mockConsulKvClient{CASModifyIndex: 3}
	handler := consulHandler{consulKvClient: consulKvClient}

	handler.handle(diffs)

	assert.Equal(t, 1, len(consulKvClient.CASs))
	expectedItem := consulKvClient.CASs[0]
	assert.Equal(t, uint64(2), expectedItem.Flags)
	assert.Equal(t, "myKey", expectedItem.Key)
	assert.Equal(t, []byte("myValue"), expectedItem.Value)
	assert.Equal(t, uint64(2), expectedItem.CreateIndex)
	assert.Equal(t, uint64(3), expectedItem.ModifyIndex)
}

func TestHandle_SetExistingItem_RetriesCASFailure(t *testing.T) {
	kv := &consul.KVPair{Flags: 2, Key: "myKey", Value: []byte("oldValue"), CreateIndex: 0, ModifyIndex: 0}
	diffs := append(differences{}, difference{
		key:            "myKey",
		sourceValue:    "myValue",
		targetValue:    "oldValue",
		targetOriginal: kv,
	})

	consulKvClient := &mockConsulKvClient{CASModifyIndex: 3, GetValue: *kv}
	handler := consulHandler{consulKvClient: consulKvClient}

	handler.handle(diffs)

	assert.Equal(t, 2, len(consulKvClient.CASs))
	expectedItem := consulKvClient.CASs[1]
	assert.Equal(t, "myKey", expectedItem.Key)
	assert.Equal(t, []byte("myValue"), expectedItem.Value)
	assert.Equal(t, uint64(3), expectedItem.ModifyIndex)
}

type mockConsulKvClient struct {
	Puts           consul.KVPairs
	CASs           consul.KVPairs
	GetValue       consul.KVPair
	CASModifyIndex uint64
}

func (client *mockConsulKvClient) Get(key string, q *consul.QueryOptions) (*consul.KVPair, *consul.QueryMeta, error) {
	client.GetValue.ModifyIndex = client.CASModifyIndex
	return &client.GetValue, nil, nil
}
func (client *mockConsulKvClient) Put(p *consul.KVPair, q *consul.WriteOptions) (*consul.WriteMeta, error) {
	client.Puts = append(client.Puts, p)
	return nil, nil
}
func (client *mockConsulKvClient) CAS(p *consul.KVPair, q *consul.WriteOptions) (bool, *consul.WriteMeta, error) {
	client.CASs = append(client.CASs, p)
	if p.ModifyIndex == client.CASModifyIndex {
		return true, nil, nil
	}
	// CAS failure
	return false, nil, nil
}
