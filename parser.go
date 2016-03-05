package consulKvMigrator

import "io/ioutil"
import "encoding/json"
import "fmt"
import "bytes"

const valuKey = "value"

func parseFromFile(filename string) []simpleKv {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	nodes := map[string]interface{}{}

	if err := json.Unmarshal(fileBytes, &nodes); err != nil {
		fmt.Println("Error parsing JSON")
		panic(err)
	}

	return parse(nodes)
}

func parse(jsonNodes map[string]interface{}) []simpleKv {
	return parseImpl(jsonNodes, []string{}, []simpleKv{})
}

func parseImpl(jsonNodes map[string]interface{}, path []string, values []simpleKv) []simpleKv {
	for key, value := range jsonNodes {
		switch value.(type) {
		case []interface{}:
			fmt.Printf("Arrays are not supported, ignoring. Key: %v", key)
		default:
			if key == valuKey {
				values = append(values, simpleKv{formatKey(path), fmt.Sprint(value)})
			} else {
				path = append(path, key)
				values = parseImpl(value.(map[string]interface{}), path, values)
				path = path[:len(path)-1]
			}
		}
	}

	return values
}

func formatKey(path []string) string {
	keyBuffer := new(bytes.Buffer)
	keyBuffer.WriteString(path[0])
	for idx := 1; idx < len(path); idx++ {
		keyBuffer.WriteString("/")
		keyBuffer.WriteString(path[idx])
	}

	return keyBuffer.String()
}

type simpleKv struct {
	key   string
	value string
}
