package contracts

import "time"

type (
	Event struct {
		ID               string    `json:"id"`
		Title            string    `json:"title"`
		Description      string    `json:"description"`
		OwnerID          string    `json:"owner_id"`
		StartDate        time.Time `json:"start_date"`
		FinishDate       time.Time `json:"fin_date"`
		NotificationTime time.Time `json:"notification_time"`
	}

	Notification struct {
		EventID    string    `json:"event_id"`
		EventTitle string    `json:"event_title"`
		EventDate  time.Time `json:"event_date"`
		UserID     string    `json:"user_id"`
	}
)
