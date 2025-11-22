package dto

type ServiceRequest struct {
	ServiceType string `json:"servicetype"`
	SlotID      int    `json:"slotid"`
}

type RescheduleServiceRequest struct {
	SlotID int `json:"slotid"`
}

type RequestProvider struct {
	AssignedTo string `json:"assignedto"`
}

type DeleteRequestMessage struct {
	UserID string `json:"userId"`
}

type Request struct {
	PK            string `dynamobdav:"PK"`
	SK            string `dynamobdav:"SK"`
	AssignedTo    string `dynamodbav:"assigned_to"`
	Date          string `dynamodbav:"date"`
	FeedbackGiven bool   `dynamodbav:"feedback_given"`
	Flat          string `dynamodbav:"flat_no"`
	ID            string `dynamodbav:"id"`
	ResidentID    string `dynamodbav:"resident_id"`
	ServiceType   string `dynamodbav:"service_type"`
	Status        string `dynamodbav:"status"`
	TimeSlot      string `dynamodbav:"time_slot"`
}
