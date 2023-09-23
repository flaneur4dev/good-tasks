package mistakes

import "errors"

var (
	ErrNoEvent     = errors.New("event not found")
	ErrDateBusy    = errors.New("this time is already captured")
	ErrCreateEvent = errors.New("failed event create")
	ErrUpdateEvent = errors.New("failed event update")
	ErrDeleteEvent = errors.New("failed event delete")
	ErrPeriod      = errors.New("incorrect start of the period")
)
