package testutils

import (
	"encoding/json"
	"io"
	"os"
)

func LoadSampleResponse[T any](endpoint string) (T, error) {
	var empty T // Equivalent to nil when we return errors
	var sample_response T
	file, err := os.Open("sample_responses/" + endpoint + ".json")
	if err != nil {
		return empty, err
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		return empty, err
	}
	err = json.Unmarshal(contents, &sample_response)
	if err != nil {
		return empty, err
	}
	return sample_response, nil
}
