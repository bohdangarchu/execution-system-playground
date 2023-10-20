package docrunner

import (
	"app/types"
	"net/http"
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
