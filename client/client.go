package prometheus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	Server *url.URL
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
		url.QueryEscape(step.String()),
	))
	if err != nil {
		return nil, err
	}

	u = c.Server.ResolveReference(u)
	r, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)

	if 400 <= r.StatusCode {
		return nil, fmt.Errorf("error response: %s", string(b))
	}

	resp := &QueryRangeResponse{}
	err = json.Unmarshal(b, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
