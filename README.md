```
 ____   _  _____ 
|  _ \ / \|_   _|
| |_) / _ \ | |  
|  __/ ___ \| |  
|_| /_/   \_\_|  
                 
```

Prometheus Alert Testing tool

Usage
=====

    pat [options] <test_yaml_file_glob>

e.g.,

    pat test/*.yaml

Sample
======

Suppose you have the following rule file that you want to be tested:

```yaml
groups:
  - name: prometheus.rules
    rules:
      - alert: HTTPRequestRateLow
        expr: http_requests{group="canary", job="app-server"} < 100
        for: 1m
        labels:
          severity: critical
```

Write a yaml file with your test cases:

```yaml
name: Test HTTP Requests too low alert
rules:
  fromFile: rules.yaml
fixtures:
  5m:
    - http_requests{job="app-server", instance="0", group="canary", severity="overwrite-me"}	75 85  95 105 105  95  85
	- http_requests{job="app-server", instance="1", group="canary", severity="overwrite-me"}	80 90 100 110 120 130 140
assertions:
  - at: 0m
    expected:
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="pending",group="canary",instance="0",job="app-server",severity="critical"} 1
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="pending",group="canary",instance="1",job="app-server",severity="critical"} 1
  - at: 5m
    expected:
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="firing",group="canary",instance="0",job="app-server",severity="critical"} 1
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="firing",group="canary",instance="1",job="app-server",severity="critical"} 1
  - at: 10m
    expected:
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="firing",group="canary",instance="0",job="app-server",severity="critical"} 1
  - at: 15m
    expected: []
  - at: 20m
    expected: []
  - at: 25m
    expected:
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="pending",group="canary",instance="0",job="app-server",severity="critical"} 1
  - at: 30m
    expected:
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="firing",group="canary",instance="0",job="app-server",severity="critical"} 1
```

    
Run the test:

```bash
$ ./pat -test.v test/test.yaml 
=== RUN   Test_HTTP_Requests_too_low_alert_0
--- PASS: Test_HTTP_Requests_too_low_alert_0 (0.00s)
=== RUN   Test_HTTP_Requests_too_low_alert_1
=== PASS: Test_HTTP_Requests_too_low_alert_1 (0.00s)
```
