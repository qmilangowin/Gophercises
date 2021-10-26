package urlshort

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type yamlData []struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func ParseYAML(data []byte) (yamlData, error) {

	var yamlData yamlData

	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return yamlData, err
	}

	return yamlData, nil

}

func BuildMap(data yamlData) map[string]string {
	pathToUrls := make(map[string]string)
	for _, i := range data {
		pathToUrls[i.Path] = i.URL
	}
	return pathToUrls
}

func toYaml(data map[string]string) error {
	var values []yamlData

	for k, v := range data {

		d := yamlData{
			{k, v},
		}
		values = append(values, d)
	}

	out, err := yaml.Marshal(&values)
	if err != nil {
		return fmt.Errorf("Could not marshal values: %v", err)
	}
	fmt.Println(string(out))
	return nil
}
