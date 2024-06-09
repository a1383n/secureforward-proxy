package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func CheckDomainAndIp(api string, domain string, ip string) (bool, error) {
	data := map[string]string{"q": domain, "ip": ip}
	dataBytes, err := json.Marshal(data)

	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", api+"/check", bytes.NewBuffer(dataBytes))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, errors.New("Check Domain Error: HTTP " + resp.Status)
	}

	return result["ok"].(bool), nil
}
