package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

// FetchJson requests the given url and then attempts to unmarshal the body as JSON into the provided struct.
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

// FetchHtmlAndFind downloads the HTML page at the given URL and runs the specified CSS selector over it to find nodes.
// The textual content of those nodes is returned.
func FetchHtmlAndFind(ctx context.Context, url string, selector string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, err
	}

	var results []string
	doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
		results = append(results, selection.Text())
	})
	return results, nil
}

// FetchHash downloads the given URL and parses the first hash out of it, assuming it's formatted in line with the
// output of sha256sum. Hashes are assumed to be hexadecimal and an error will be returned if this is not the case.
func FetchHash(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	hash := strings.ToLower(strings.SplitN(string(b), " ", 2)[0])
	for i := range hash {
		if (hash[i] < 'a' || hash[i] > 'f') && (hash[i] < '0' || hash[i] > '9') {
			return "", fmt.Errorf("invalid has found at address: %s", hash)
		}
	}
	return hash, nil
}
