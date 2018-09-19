package log

import (
	"fmt"
)

// SeriesName specifies the key used to define which bucket/table to send values to
// thius is used by adapters to extract the table name from values supplied.
const SeriesName = "stats_series_name"

// StatsLog conforms to the ValuesLogger interface
// It simply outputs the values to stdout, rather than logging to a specific service
// see the adapters folder for ValueLoggers which connect to time series databases.
type StatsLog struct {
}

// Values prints to the ValuesLog which typically emits stats to a time series database.
// This example logger prints values to stdout instead.
func (l *StatsLog) Values(values map[string]interface{}) {
	fmt.Printf("Values logged:%+s", values)
}
