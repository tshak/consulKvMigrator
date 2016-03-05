package consulKvMigrator

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func TestParse(t *testing.T) {
	migrations := parseFromFile("parser_test.json")

	migrationsMap := map[string]string{}
	for idx := range migrations {
		migrationsMap[migrations[idx].key] = migrations[idx].value
	}

	assert.Equal(t, 4, len(migrations))
	assert.Equal(t, "value1", migrationsMap["level1/key1"], fmt.Sprintf("%s", migrationsMap))
	assert.Equal(t, "value2", migrationsMap["level1/level2/level3/key2"])
	assert.Equal(t, "value3", migrationsMap["level1/level2/key3"])
	assert.Equal(t, "0.42", migrationsMap["key4"])
}
