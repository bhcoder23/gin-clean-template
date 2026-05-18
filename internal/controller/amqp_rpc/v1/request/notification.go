package request

// ListNotifications -.
type ListNotifications struct {
	UnreadOnly bool `json:"unread_only"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
}

// MarkNotificationRead -.
type MarkNotificationRead struct {
	ID string `json:"id" validate:"required"`
}
