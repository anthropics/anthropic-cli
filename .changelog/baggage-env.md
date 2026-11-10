feat(client): read OpenTelemetry baggage from the BAGGAGE environment variable

While tracing is on, API requests send the value of `BAGGAGE` in the `baggage` header, the same way they continue the trace in `TRACEPARENT`.
