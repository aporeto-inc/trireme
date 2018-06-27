package statscollector

import (
	"sync"

	"go.aporeto.io/trireme-lib/collector"
)

// NewCollector provides a new collector interface
func NewCollector() Collector {
	return &collectorImpl{
		Flows:          map[string]*collector.FlowRecord{},
		Users:          map[string]*collector.UserRecord{},
		ProcessedUsers: map[string]bool{},
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
	Flows          map[string]*collector.FlowRecord
	ProcessedUsers map[string]bool
	Users          map[string]*collector.UserRecord
	sync.Mutex
}
