package other

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

// GetRequest executes a GET request
func GetRequest(url string) ([]byte, error) {

	client := http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	reqWithDeadline := req.WithContext(ctx)
	response, err := client.Do(reqWithDeadline)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	return data, err

}
