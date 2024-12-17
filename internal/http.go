package internal

import (
	"context"
	"encoding/json"
	"net/http"
)

func FetchJson(ctx context.Context, url string, i interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(i)
}
