package models

import (
	"time"
)

type (
	Build struct {
		id string
		last_run string
		run_start_time time.Time
		log string
		status string
	}
)