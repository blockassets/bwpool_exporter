package main

import (
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/blockassets/bwpool_exporter/bwpool"
	"github.com/blockassets/prometheus_helper"
	"github.com/prometheus/client_golang/prometheus"
)

//
var (
	namespace    = "bwpool"
	idLabelNames = []string{"id"}
)

// Collector interface
type Exporter struct {
	client      *bwpool.BWClient
	ConstLabels prometheus.Labels
	Gauges      prometheus_helper.GaugeMapMap
	GaugeVecs   prometheus_helper.GaugeVecMapMap
	sync.Mutex
}

//
func NewExporter(poolConfig *bwpool.PoolConfig, timeout time.Duration) *Exporter {
	constLabels := prometheus.Labels{"key": poolConfig.PublicKey[:8]}
	structFieldMap := prometheus_helper.NewStructFieldMap(bwpool.PoolData{})

	exporter := &Exporter{
		client:      bwpool.NewClient(poolConfig, timeout),
		ConstLabels: constLabels,
		Gauges:      prometheus_helper.NewGaugeMapMap(structFieldMap, namespace, constLabels),
		GaugeVecs:   make(prometheus_helper.GaugeVecMapMap),
	}

	return exporter
}

//
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	prometheus_helper.DescribeGaugeMapMap(e.Gauges, ch)
	prometheus_helper.DescribeGaugeVecMapMap(e.GaugeVecs, ch)
}

//
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Prevents multiple concurrent calls
	e.Lock()
	defer e.Unlock()

	poolData, err := e.client.FetchWorkers()
	if err != nil {
		log.Println(err)
		return
	}

	poolDataMap := prometheus_helper.NewStructFieldMap(*poolData)

	for key, value := range poolDataMap {
		val := reflect.ValueOf(value)
		// 'Workers' is a special case as a GaugeVec
		if val.Kind() == reflect.Map {
			for _, key := range val.MapKeys() {
				name := key.Interface().(string)
				worker := val.MapIndex(key).Interface()
				prometheus_helper.CollectGaugeVecs(name, worker, e.GaugeVecs, namespace, e.ConstLabels, idLabelNames, prometheus.Labels{idLabelNames[0]: name})
			}
		} else {
			meta := prometheus_helper.StructMeta{}
			prometheus_helper.MakeStructMeta(value, &meta)
			prometheus_helper.SetValuesOnGauges(meta, namespace, e.Gauges[key])
		}
	}

	prometheus_helper.CollectGaugeMapMap(e.Gauges, ch)
	prometheus_helper.CollectGaugeVecMapMap(e.GaugeVecs, ch)
}
