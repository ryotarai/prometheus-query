package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	prometheus "github.com/ryotarai/prometheus-query/client"
	"github.com/ymotongpoo/datemaki"
)

type options struct {
	format string
	server string
	query  string
	start  string
	end    string
	step   string
}

func main() {
	options := parseFlags()
	err := validateOptions(options)
	if err != nil {
		onError(err)
	}

	start, err := datemaki.Parse(options.start)
	if err != nil {
		onError(err)
	}

	end, err := datemaki.Parse(options.end)
	if err != nil {
		onError(err)
	}

	step, err := time.ParseDuration(options.step)
	if err != nil {
		onError(err)
	}

	client, err := prometheus.NewClient(options.server)
	if err != nil {
		onError(err)
	}

	resp, err := client.QueryRange(options.query, start, end, step)
	if err != nil {
		onError(err)
	}

	err = printResp(resp, options.format)
	if err != nil {
		onError(err)
	}
}

func onError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func parseFlags() options {
	format := flag.String("format", "json", "Format (available formats are json, tsv and csv)")
	server := flag.String("server", os.Getenv("PROMETHEUS_SERVER"), "Prometheus server URL like 'https://prometheus.example.com' (can be set by PROMETHEUS_SERVER environment variable)")
	query := flag.String("query", "", "Query")
	start := flag.String("start", "1 hour ago", "Start time")
	end := flag.String("end", "now", "End time")
	step := flag.String("step", "15s", "Step")

	flag.Parse()

	return options{
		format: *format,
		server: *server,
		query:  *query,
		start:  *start,
		end:    *end,
		step:   *step,
	}
}

func validateOptions(options options) error {
	if options.server == "" {
		return errors.New("-server is mandatory")
	}
	if options.query == "" {
		return errors.New("-query is mandatory")
	}

	return nil
}
