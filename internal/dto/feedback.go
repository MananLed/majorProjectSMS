package dto

type Feedback struct {
	Rating    int32  `json:"rating"`
	Content   string `json:"content"`
	RequestID string `json:"requestid"`
}
