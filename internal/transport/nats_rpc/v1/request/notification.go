package request

// ListNotificationsReq -.
type ListNotificationsReq struct {
	UnreadOnly bool `json:"unread_only"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
}

// MarkNotificationReadReq -.
type MarkNotificationReadReq struct {
	ID string `json:"id" validate:"required"`
}
