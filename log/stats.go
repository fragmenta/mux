package log

import (
	"fmt"
)

// Emit stats to a time series database, conform to ValueLogger interface
// need adapters for different db backends, perhaps just start with influx though

//StatsLog conforms to the ValuesLogger interface
type StatsLog struct {
}

// Values prints to the ValuesLog which typically emits stats to a time series database.
func (l *StatsLog) Values(values map[string]interface{}) {

	// for now just print to stdout
	fmt.Printf("VALUES:%s", values)
	//l.Values(values)
}
