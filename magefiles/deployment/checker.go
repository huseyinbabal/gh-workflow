package deployment

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Status struct {
	Version string `json:"version"`
	Ready   bool   `json:"ready"`
}

func WaitFor(version string) error {
	deadline := time.Now().Add(100 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(2 * time.Second)
		status, err := getDeploymentStatus()
		if err != nil {
			fmt.Println("failed to get status.", err)
			continue
		}
		if status.Version == version {
			return nil
		}
		fmt.Println("version mismatch. actual:", status.Version, ", expected: ", version)
	}
	return errors.New("deployment verification timeout")
}

func getDeploymentStatus() (*Status, error) {
	resp, err := http.Get("https://webhook.site/8fc86320-5be4-46e7-a454-e45f45f297b1")
	if err != nil {
		return nil, fmt.Errorf("while getting deployment status %w", err)
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading response body %w", err)
	}

	var st Status
	err = json.Unmarshal(bytes, &st)
	if err != nil {
		return nil, fmt.Errorf("while unmarshalling struct %w", err)
	}
	return &st, nil
}
