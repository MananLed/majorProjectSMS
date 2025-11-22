package dto

type Feedback struct {
	Rating    int32  `json:"rating"`
	Content   string `json:"content"`
	RequestID string `json:"requestid"`
}

type FeedbackDDB struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	AssignedTo  string `dynamodbav:"assigned_to"`
	Content     string `dynamodbav:"content"`
	Date        string `dynamodbav:"date"`
	Flat        string `dynamodbav:"flat_no"`
	ID          string `dynamodbav:"id"`
	Rating      int32  `dynamodbav:"rating"`
	RequestID   string `dynamodbav:"request_id"`
	ResidentID  string `dynamodbav:"resident_id"`
	ServiceType string `dynamodbav:"service_type"`
	TimeSlot    string `dynamodbav:"time_slot"`
	UserName    string `dynamodbav:"username"`
}
