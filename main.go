package main

import (
    "flag"
    "bytes"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/log"
    "github.com/prometheus/common/expfmt"
    "github.com/gjflsl/nginx_exporter/ngx_status_exporter"
    req_status_exporter "github.com/gjflsl/nginx_exporter/req_status_exporter"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
    namespace = "nginx" // For Prometheus metrics.
)

var (
    VERSION = "1.0.1"

    listeningAddress = flag.String("telemetry.address", ":9113", "Address on which to expose metrics.")
    metricsEndpoint  = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
)

func main() {
    flag.Parse()

    log.Printf("Starting Server: %s", *listeningAddress)
    log.Printf("Nginx Exporter %s started.", VERSION)
    http.Handle(*metricsEndpoint, promhttp.Handler())
    http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
        params := r.URL.Query()

        nginxStatus := params.Get("module")
        nginxScrapeURI := params.Get("target")

        if nginxStatus == "" {
            nginxStatus = "ngx_status"
        }

        if nginxScrapeURI == "" {
            http.Error(w, "Url parameter is missing", 400)
            return
        }

        if nginxStatus == "ngx_status" || nginxStatus == "" {

            exp := ngx_status_exporter.NewNgxStatusExporter(nginxScrapeURI, true)
            registry := prometheus.NewRegistry()
            registry.Register(exp)
            mfs, err := registry.Gather()

            if err != nil {
                log.Fatal(err)
            }
            var buf bytes.Buffer
            for _, mf := range mfs {
                if _, err := expfmt.MetricFamilyToText(&buf, mf); err != nil {
                    log.Fatal(err)
                }
            }
            w.Write(buf.Bytes())
        } else if nginxStatus == "req_status" {
            exp := req_status_exporter.NewReqStatusExporter(nginxScrapeURI, true)
            registry := prometheus.NewRegistry()
            registry.Register(exp)
            mfs, err := registry.Gather()

            if err != nil {
                log.Fatal(err)
            }
            var buf bytes.Buffer
            for _, mf := range mfs {
                if _, err := expfmt.MetricFamilyToText(&buf, mf); err != nil {
                    log.Fatal(err)
                }
            }
            w.Write(buf.Bytes())
        } else {
            http.Error(w, "Type parameter is error", 400)
        }

    })
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`<html>what the fuck!</html>`))
    })

    log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
