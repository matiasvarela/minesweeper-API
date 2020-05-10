package clock

import "time"

//go:generate mockgen -source=clock.go -destination=../../mock/clock.go -package=mock

type Clock interface {
	Now() time.Time
}

type clock struct{}

func New() Clock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now().UTC()
}
