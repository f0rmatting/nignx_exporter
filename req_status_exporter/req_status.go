package rt_status_exporter

import (
    "crypto/tls"
    "net/http"
    "strings"
    "sync"
    "strconv"
    "fmt"
    "io/ioutil"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/log"
)

const (
    namespace = "nginx" // For Prometheus metrics.
)

// Exporter collects nginx stats from the given URI and exports them using
// the prometheus metrics package.
type Exporter struct {
    URI    string
    mutex  sync.RWMutex
    client *http.Client

    scrapeFailures prometheus.Counter
    nginxReqStatus *prometheus.Desc
}

// NewExporter returns an initialized Exporter.
func NewReqStatusExporter(uri string, insecure bool) *Exporter {
    return &Exporter{
        URI: uri,
        scrapeFailures: prometheus.NewCounter(prometheus.CounterOpts{
            Namespace: namespace,
            Name:      "exporter_scrape_failures_total",
            Help:      "Number of errors while scraping nginx.",
        }),
        nginxReqStatus: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "rs_status"),
            "Number of req status by tengine",
            []string{"domain", "key_type"},
            nil,
        ),
        client: &http.Client{
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
            },
        },
    }
}

// Describe describes all the metrics ever exported by the nginx exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
    e.scrapeFailures.Describe(ch)
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
    reqStatuKeys := [] string{"kv", "bytes_in", "bytes_out", "conn_total", "req_total", "http_2xx", "http_3xx", "http_4xx", "http_5xx", "http_other_status", "rt", "ups_req", "ups_rt", "ups_tries", "http_200", "http_206", "http_302", "http_304", "http_403", "http_404", "http_416", "http_499", "http_500", "http_502", "http_503", "http_504", "http_508", "http_other_detail_status", "http_ups_4xx", "http_ups_5xx"}
    resp, err := e.client.Get(e.URI)
    if err != nil {
        return fmt.Errorf("Error scraping nginx: %v", err)
    }

    data, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 400 {
        if err != nil {
            data = []byte(err.Error())
        }
        return fmt.Errorf("Status %s (%d): %s", resp.Status, resp.StatusCode, data)
    }

    // Parsing results
    lines := strings.Split(string(data), "\n")

    countResult := make(map[string][]int)
    for _, line := range lines {
        reqStatusValues := strings.Split(line, ",")
        if len(reqStatusValues) >= 30 {
            keyFixNum := len(reqStatusValues) - 29
            kvArr := reqStatusValues[:keyFixNum]
            domain := strings.Join(kvArr, ",")
            countResult[domain] = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
            for index, value := range reqStatusValues[keyFixNum:] {
                v, err := strconv.Atoi(value)
                if err != nil {
                    return err
                }
                countResult[domain][index] += v
            }
        }

    }

    for domain, reqDataList := range countResult {
        for index, value := range reqDataList {
            ch <- prometheus.MustNewConstMetric(e.nginxReqStatus, prometheus.CounterValue, float64(value), domain, reqStatuKeys[index+1])
        }
    }

    return nil
}

// Collect fetches the stats from configured nginx location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
    e.mutex.Lock() // To protect metrics from concurrent collects.
    defer e.mutex.Unlock()
    if err := e.collect(ch); err != nil {
        log.Printf("Error scraping nginx: %s", err)
        e.scrapeFailures.Inc()
        e.scrapeFailures.Collect(ch)
    }
    return
}
