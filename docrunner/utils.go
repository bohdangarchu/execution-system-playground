package docrunner

import (
	"app/types"
	"net/http"
	"time"
)

func CheckContainerHealth(conatiner *types.DockerContainer) bool {
	url := "http://localhost:" + conatiner.Port + "/health"
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func WaitUntilAvailable(container *types.DockerContainer) {
	for {
		if CheckContainerHealth(container) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
