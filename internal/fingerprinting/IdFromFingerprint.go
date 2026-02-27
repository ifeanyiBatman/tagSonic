package fingerprinting

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AcousticIDResponse struct {
	Status string `json:"status"`
	Error  struct {
		Message string `json:"message"`
	} `json:"error"`
	Results []struct {
		ID    string  `json:"id"`
		Score float64 `json:"score"`
	} `json:"results"`
}

func GetIdFromFingerprint(fingerprint string, duration float64, APIKey string) (string, error) {
	lookupURL := "https://api.acoustid.org/v2/lookup"

	req, err := http.NewRequest("GET", lookupURL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("client", APIKey)
	q.Add("fingerprint", fingerprint)
	q.Add("duration", fmt.Sprintf("%d", int(duration)))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var acousticIDResponse AcousticIDResponse
	if err := json.NewDecoder(res.Body).Decode(&acousticIDResponse); err != nil {
		return "", err
	}

	if acousticIDResponse.Status == "error" {
		return "", fmt.Errorf("AcoustID API error: %s", acousticIDResponse.Error.Message)
	}

	if len(acousticIDResponse.Results) == 0 {
		return "", fmt.Errorf("no results found for fingerprint")
	}

	return acousticIDResponse.Results[0].ID, nil

}
