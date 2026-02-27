package acoustid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type acoustIDResponse struct {
	Status string `json:"status"`
	Error  struct {
		Message string `json:"message"`
	} `json:"error"`
	Results []struct {
		ID         string  `json:"id"`
		Score      float64 `json:"score"`
		Recordings []struct {
			Title   string `json:"title"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"recordings"`
	} `json:"results"`
}

type SongMetadata struct {
	ID     string
	Title  string
	Artist string
}

func LookupMetadata(fingerprint string, duration float64, apiKey string) (*SongMetadata, error) {
	lookupURL := "https://api.acoustid.org/v2/lookup"

	req, err := http.NewRequest("GET", lookupURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("client", apiKey)
	q.Add("fingerprint", fingerprint)
	q.Add("duration", fmt.Sprintf("%d", int(duration)))
	q.Add("meta", "recordings")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response acoustIDResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Status == "error" {
		return nil, fmt.Errorf("AcoustID API error: %s", response.Error.Message)
	}

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results found for fingerprint")
	}

	result := response.Results[0]

	meta := &SongMetadata{
		ID: result.ID,
	}

	if len(result.Recordings) > 0 {
		meta.Title = result.Recordings[0].Title
		if len(result.Recordings[0].Artists) > 0 {
			meta.Artist = result.Recordings[0].Artists[0].Name
		}
	}

	return meta, nil
}
