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
	Version string
	Ready   bool
}

func WaitFor(url, version string) error {
	deadline := time.Now().Add(100 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(2 * time.Second)
		status, err := getDeploymentStatus(url)
		if err != nil {
			fmt.Println("failed to get status.", err)
			continue
		}
		if status.Version == version && status.Ready {
			fmt.Println("deployment is ready.")
			return nil
		}
		fmt.Println("deployment is not ready, actual version:", status.Version, ", expected version: ", version, ", ready: ", status.Ready)
	}
	return errors.New("deployment verification timeout")
}

func getDeploymentStatus(url string) (*Status, error) {
	resp, err := http.Get(url)
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
