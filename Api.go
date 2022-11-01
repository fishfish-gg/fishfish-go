package fishfish

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

func (client *FishFishClient) getAPIUrl(path string) string {
	return fmt.Sprintf("%s%s", client.url, path)
}

func (client *FishFishClient) authenticatedRequest(path, requestType, body string) (*http.Response, error) {
	req, _ := http.NewRequest(strings.ToUpper(requestType), client.getAPIUrl(path), bytes.NewBufferString(body))
	req.Header.Set("Authorization", client.getSessionToken())
	resp, err := client.httpClient.Do(req)

	return resp, err
}
