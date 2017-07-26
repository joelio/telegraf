package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var (
	BatchSize         int
	MaxBackLog        int
	BatchTimeInterval int
	SpanCount         int
	ZipkinServerHost  string
)

const usage = `./stress_test_write -batch_size=<batch_size> -max_backlog=<max_span_buffer_backlog> -batch_interval=<batch_interval_in_seconds> -span_count<number_of_spans_to_write> -zipkin_host=<zipkin_service_hostname>`

func init() {
	flag.IntVar(&BatchSize, "batch_size", 10000, usage)
	flag.IntVar(&MaxBackLog, "max_backlog", 100000, usage)
	flag.IntVar(&BatchTimeInterval, "batch_interval", 1, usage)
	flag.IntVar(&SpanCount, "span_count", 100000, usage)
	flag.StringVar(&ZipkinServerHost, "zipkin_host", "localhost", usage)
}

func main() {
	flag.Parse()
	var hostname = fmt.Sprintf("http://%s:9411/api/v1/spans", ZipkinServerHost)
	collector, err := zipkin.NewHTTPCollector(
		hostname,
		zipkin.HTTPBatchSize(BatchSize),
		zipkin.HTTPMaxBacklog(MaxBackLog),
		zipkin.HTTPBatchInterval(time.Duration(BatchTimeInterval)*time.Second))
	defer collector.Close()
	if err != nil {
		log.Fatalf("Error intializing zipkin http collector: %v\n", err)
	}

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "127.0.0.1:0", "trivial"))

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	log.Printf("Writing %d spans to zipkin server at %s\n", SpanCount, hostname)
	for i := 0; i < SpanCount; i++ {
		parent := tracer.StartSpan("Parent")
		parent.LogEvent(fmt.Sprintf("Trace%d", i))
		parent.Finish()
	}
	log.Println("Done. Flushing remaining spans...")
}
