package fishfish

import (
	"fmt"
)

func (client *Client) getAPIUrl(path string) string {
	return fmt.Sprintf("%s%s", client.url, path)
}
