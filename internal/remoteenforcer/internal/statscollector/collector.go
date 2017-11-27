package statscollector

import (
	"sync"

	"github.com/aporeto-inc/trireme-lib/collector"
)

// NewCollector provides a new collector interface
func NewCollector() Collector {
	return &collectorImpl{
		Flows: map[string]*collector.FlowRecord{},
	}
}

// collectorImpl : This object is a stash implements two interfaces.
//
//  collector.EventCollector - so datapath can report flow events
//  CollectorReader - so components can extract information out of this stash
//
// It has a flow entries cache which contains unique flows that are reported
// back to the controller/launcher process
type collectorImpl struct {
	Flows map[string]*collector.FlowRecord
	sync.Mutex
}
