package chartmuseum

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"io/ioutil"
)

// DownloadFile downloads a file from ChartMuseum
func (client *Client) DownloadFile(filePath string) ([]byte, error) {
	u, err := url.Parse(client.opts.url)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(client.opts.contextPath, strings.TrimPrefix(u.Path, client.opts.contextPath), filePath)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	if client.opts.accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.opts.accessToken))
	} else if client.opts.username != "" && client.opts.password != "" {
		req.SetBasicAuth(client.opts.username, client.opts.password)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}
