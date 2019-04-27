package main

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	loginFailedTotal     = stats.Int64("login_failed_total", "Total number of succeeded logins", "1")
	loginFailedTotalView = &view.View{
		Name:        "login_failed_total",
		Measure:     loginFailedTotal,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(loginFailedTotalView)
}
