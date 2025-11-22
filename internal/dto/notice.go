package dto

type NoticeRequest struct {
	Content string `json:"content"`
}

type Notice struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	ID          string `dynamodbav:"id"`
	Content     string `dynamodbav:"content"`
	Date_Issued string `dynamodbav:"date_issued"`
	Month       string `dynamodbav:"month"`
	Year        int    `dynamodbav:"year"`
}
