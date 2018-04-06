package crawler

import (
	"sync"
)

// InstrumentationMem stores the metric data from the application in memory
type InstrumentationMem struct {
	counts    map[string]int64
	gauges    map[string]int64
	histogram map[string]hist

	countMU sync.Mutex
	gaugeMU sync.Mutex
	histMU  sync.Mutex
}

type hist struct {
	values map[string]int64
}

func (h hist) add(value string) {
	h.values[value]++
}

// NewInstrumentationMem builds a basic in-memory instrumentation
func NewInstrumentationMem() Instrument {
	return &InstrumentationMem{
		counts:    make(map[string]int64),
		gauges:    make(map[string]int64),
		histogram: make(map[string]hist),
	}
}

// Count a metric by increasing its value
func (imem *InstrumentationMem) Count(metric string) {
	imem.countMU.Lock()
	defer imem.countMU.Unlock()

	imem.counts[metric]++
}

// Gauge changes a metric by a given integer amount
func (imem *InstrumentationMem) Gauge(metric string, value int64) {
	imem.gaugeMU.Lock()
	defer imem.gaugeMU.Unlock()
	imem.gauges[metric] += value
}

// Histogram tores received values for a given metric in a countable way
func (imem *InstrumentationMem) Histogram(metric string, value string) {
	imem.histMU.Lock()
	defer imem.histMU.Unlock()

	_, ok := imem.histogram[metric]
	if !ok {
		imem.histogram[metric] = hist{values: make(map[string]int64)}
	}

	imem.histogram[metric].add(value)
}
