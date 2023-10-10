package contracts

import "time"

type (
	Event struct {
		ID              string    `json:"id"`
		Title           string    `json:"title"`
		Description     string    `json:"description"`
		OwnerID         string    `json:"owner_id"`
		StartDate       time.Time `json:"start_date"`
		FinishDate      time.Time `json:"finish_date"`
		NotificationDay time.Time `json:"notification_day"`
	}

	Notification struct {
		UserID     string    `json:"user_id"`
		EventID    string    `json:"event_id"`
		EventTitle string    `json:"event_title"`
		EventDate  time.Time `json:"event_date"`
	}

	NotificationMessage struct {
		ID              string `json:"id"`
		ContentType     string `json:"content_type"`
		ContentEncoding string `json:"content_encoding"`
		Body            []byte `json:"body"`
	}
)
