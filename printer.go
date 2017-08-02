package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	prometheus "github.com/ryotarai/prometheus-query/client"
)

func printResp(resp *prometheus.QueryRangeResponse, format string) error {
	switch format {
	case "tsv":
		return printRespXSV(resp, "\t")
	case "csv":
		return printRespXSV(resp, ",")
	case "json":
		return printRespJSON(resp)
	}

	return fmt.Errorf("unknown format: %s", format)
}

func printRespJSON(resp *prometheus.QueryRangeResponse) error {
	type valueEntry struct {
		Metric map[string]string `json:"metric"`
		Value  float64           `json:"value"`
	}
	type timeEntry struct {
		Time   int64         `json:"time"`
		Values []*valueEntry `json:"values"`
	}
	entryByTime := map[int64]*timeEntry{}

	for _, r := range resp.Data.Result {
		for _, v := range r.Values {
			t := v.Time()
			u := t.Unix()
			e, ok := entryByTime[u]
			if !ok {
				e = &timeEntry{
					Time:   u,
					Values: []*valueEntry{},
				}
				entryByTime[u] = e
			}

			val, err := v.Value()
			if err != nil {
				return err
			}
			e.Values = append(e.Values, &valueEntry{
				Metric: r.Metric,
				Value:  val,
			})
		}
	}

	s := make([]*timeEntry, len(entryByTime))
	i := 0
	for _, e := range entryByTime {
		s[i] = e
		i++
	}

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func printRespXSV(resp *prometheus.QueryRangeResponse, delimiter string) error {
	type valueByMetric map[string]float64

	valuesByTime := map[time.Time]valueByMetric{}
	metrics := []string{}

	for _, r := range resp.Data.Result {
		metric := stringMapToString(r.Metric, "|")
		for _, v := range r.Values {
			t := v.Time()
			d, ok := valuesByTime[t]
			if !ok {
				d = valueByMetric{}
				valuesByTime[t] = d
			}
			var err error
			d[metric], err = v.Value()
			if err != nil {
				return err
			}
		}

		found := false
		for _, m := range metrics {
			if m == metric {
				found = true
			}
		}
		if !found {
			metrics = append(metrics, metric)
		}
	}

	type st struct {
		time time.Time
		v    valueByMetric
	}
	slice := make([]st, len(valuesByTime))
	i := 0
	for t, v := range valuesByTime {
		slice[i] = st{t, v}
		i++
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].time.Before(slice[j].time)
	})

	// header
	fmt.Printf("time%s%s\n", delimiter, strings.Join(metrics, delimiter))

	// print rows
	for _, s := range slice {
		values := make([]string, len(metrics))
		for i, m := range metrics {
			if v, ok := s.v[m]; ok {
				values[i] = fmt.Sprintf("%f", v)
			} else {
				values[i] = ""
			}
		}
		fmt.Printf("%d%s%s\n", s.time.Unix(), delimiter, strings.Join(values, delimiter))
	}

	return nil
}
