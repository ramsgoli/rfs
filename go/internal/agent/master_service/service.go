package masterservice

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type MasterService interface {
	RegisterWithMaster(masterHostName string, masterPort int64, agentIpAddress string, agentPort int64) error
}

func NewMasterService(httpClient HttpClient) MasterService {
	return &masterService{
		httpClient: httpClient,
	}
}

type masterService struct {
	httpClient HttpClient
}

func (m *masterService) RegisterWithMaster(masterHostName string, masterPort int64, agentIpAddress string, agentPort int64) error {
	body := []byte(fmt.Sprintf(`{
		"port": %d,
		"ip_address": "%s"
	}`, agentPort, agentIpAddress))

	resp, err := m.httpClient.Post(fmt.Sprintf("http://%s:%d/node/register", masterHostName, masterPort), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error registering with master: %v", err)
	}

	var response registerResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if !response.Success {
		return fmt.Errorf("failed to register with master")
	}
	return nil
}

type registerResponse struct {
	Success bool `json:"success"`
}
