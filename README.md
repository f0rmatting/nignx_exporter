# Nginx Exporter for Prometheus

Fork from https://github.com/discordianfish/nginx_exporter
ADD tengine req_status support.

To run it:

```
./nginx_exporter [flags]
```

Help on flags:

```
./nginx_exporter --help
```

Get nginx status

```
curl http://127.0.0.1:9113/probe?module=ngx_status&target=http%3A%2F%2F172.19.0.116%2Fstatus

# HELP nginx_connections_current Number of connections currently processed by nginx
# TYPE nginx_connections_current gauge
nginx_connections_current{state="active"} 1077
nginx_connections_current{state="reading"} 0
nginx_connections_current{state="waiting"} 1072
nginx_connections_current{state="writing"} 5
# HELP nginx_connections_processed_total Number of connections processed by nginx
# TYPE nginx_connections_processed_total counter
nginx_connections_processed_total{stage="accepted"} 2.50746094e+08
nginx_connections_processed_total{stage="handled"} 2.50746094e+08
nginx_connections_processed_total{stage="request_time"} 8.584747354e+09
nginx_connections_processed_total{stage="requests"} 1.24481088e+08

```

Get tengine req_status

```
curl http://127.0.0.1:9113/probe?module=req_status&target=http%3A%2F%2F172.20.4.14%2Freq_status

nginx_rs_status{domain="172.20.4.14",key_type="bytes_in"} 2992
nginx_rs_status{domain="172.20.4.14",key_type="bytes_out"} 12254
nginx_rs_status{domain="172.20.4.14",key_type="conn_total"} 22
nginx_rs_status{domain="172.20.4.14",key_type="http_200"} 22
nginx_rs_status{domain="172.20.4.14",key_type="http_206"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_2xx"} 22
nginx_rs_status{domain="172.20.4.14",key_type="http_302"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_304"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_3xx"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_403"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_404"} 1
nginx_rs_status{domain="172.20.4.14",key_type="http_416"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_499"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_4xx"} 1
nginx_rs_status{domain="172.20.4.14",key_type="http_500"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_502"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_503"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_504"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_508"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_5xx"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_other_detail_status"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_other_status"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_ups_4xx"} 0
nginx_rs_status{domain="172.20.4.14",key_type="http_ups_5xx"} 0
nginx_rs_status{domain="172.20.4.14",key_type="req_total"} 23
nginx_rs_status{domain="172.20.4.14",key_type="rt"} 0
nginx_rs_status{domain="172.20.4.14",key_type="ups_req"} 0
nginx_rs_status{domain="172.20.4.14",key_type="ups_rt"} 0
nginx_rs_status{domain="172.20.4.14",key_type="ups_tries"} 0
```

Prometheus congfig

```
  - job_name: 'nginx_status_exporter'
    metrics_path: /probe
    targets:
      ["172.52.33.11", "172.11.33.25"]
    params:
      module: [ngx_status]
    relabel_configs:
      - source_labels: [__address__]
        regex: (.*)
        target_label: __param_target
        replacement: ${1}
      - source_labels: [__param_target]
        regex: (.*)
        target_label: check_addr
        replacement: ${1}
      - source_labels: []
        regex: .*
        target_label: __address__
        replacement: 127.0.0.1:9113 # You nginx_exporter
        
  - job_name: 'req_status_exporter'
      metrics_path: /probe
      targets:
        ["172.52.33.11", "172.11.33.25"]
      params:
        module: [req_status]
      relabel_configs:
        - source_labels: [__address__]
          regex: (.*)
          target_label: __param_target
          replacement: ${1}
        - source_labels: [__param_target]
          regex: (.*)
          target_label: check_addr
          replacement: ${1}
        - source_labels: []
          regex: .*
          target_label: __address__
          replacement: 127.0.0.1:9113 # You nginx_exporter      
```

