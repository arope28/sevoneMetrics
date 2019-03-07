# sevoneMetrics

- *sevone_metrics.go* is WIP - it's supposed to have concurrency added, but it's not there yet
- *sevone_metrics_concurrency.go* adds concurrency. Without concurrency (*sevone_metrics.go*), the script takes ~4m51s to pull back 4232 metrics for 106 devices. WITH concurrency (this script), the same 4232 metrics for 106 devices takes ~9.13 seconds.
