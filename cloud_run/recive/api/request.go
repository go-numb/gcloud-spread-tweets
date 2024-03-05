package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (p *Client) RequestForPost(v url.Values) error {
	u, err := url.Parse(p.POSTAPIURL)
	if err != nil {
		return err
	}

	u.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	statusCode, err := _request(req, nil)
	if err != nil {
		return err
	} else if statusCode != http.StatusOK {
		return fmt.Errorf("failed to request to post process: %d", statusCode)
	}

	return nil
}

func _request(req *http.Request, binder any) (int, error) {
	hc := http.DefaultClient
	res, err := hc.Do(req)
	if err != nil {
		return res.StatusCode, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, err
	}

	if binder != nil {
		if err := json.Unmarshal(b, binder); err != nil {
			return res.StatusCode, err
		}
	}

	return res.StatusCode, nil
}
