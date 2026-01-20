package vehicle

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func FetchCarSpec(makeName, modelName string) (*CarSpecResponse, error) {
	url := fmt.Sprintf(
		"https://cars-by-api-ninjas.p.rapidapi.com/v1/cars?make=%s&model=%s",
		makeName, modelName,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-RapidAPI-Key", "<API_KEY>")
	req.Header.Set("X-RapidAPI-Host", "cars-by-api-ninjas.p.rapidapi.com")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cars API returned status %d", res.StatusCode)
	}

	var cars []CarSpecResponse
	if err := json.NewDecoder(res.Body).Decode(&cars); err != nil {
		return nil, err
	}

	if len(cars) == 0 {
		return nil, errors.New("vehicle not found in external API")
	}

	return &cars[0], nil
}
