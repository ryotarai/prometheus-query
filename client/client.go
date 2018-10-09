package prometheus

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	Server  *url.URL
	Auth    string
	hasAuth bool
}

func NewClient(addr string) (*Client, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		Server: u,
	}, nil
}

func NewAuthClient(addr, user, password string) (*Client, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		Server:  u,
		Auth:    basicAuth(user, password),
		hasAuth: true,
	}, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type QueryRangeResponse struct {
	Status string                  `json:"status"`
	Data   *QueryRangeResponseData `json:"data"`
}

type QueryRangeResponseData struct {
	Result []*QueryRangeResponseResult `json:"result"`
}

type QueryRangeResponseResult struct {
	Metric map[string]string          `json:"metric"`
	Values []*QueryRangeResponseValue `json:"values"`
}

type QueryRangeResponseValue []interface{}

func (v *QueryRangeResponseValue) Time() time.Time {
	t := (*v)[0].(float64)
	return time.Unix(int64(t), 0)
}

func (v *QueryRangeResponseValue) Value() (float64, error) {
	s := (*v)[1].(string)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

func (c *Client) QueryRange(query string, start time.Time, end time.Time, step time.Duration) (*QueryRangeResponse, error) {
	u, err := url.Parse(fmt.Sprintf("./api/v1/query_range?query=%s&start=%s&end=%s&step=%s",
		url.QueryEscape(query),
		url.QueryEscape(fmt.Sprintf("%d", start.Unix())),
		url.QueryEscape(fmt.Sprintf("%d", end.Unix())),
		url.QueryEscape(fmt.Sprintf("%ds", int(step.Seconds()))),
	))
	if err != nil {
		return nil, err
	}

	u = c.Server.ResolveReference(u)

	req, err := http.NewRequest("GET", u.String(), nil)

	if c.hasAuth {
		req.Header.Add("Authorization", "Basic "+c.Auth)
	}

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if err == io.EOF {
			return &QueryRangeResponse{}, nil
		}
		return nil, err
	}

	if 400 <= res.StatusCode {
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	resp := &QueryRangeResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
