package dto

type ServiceRequest struct {
	ServiceType string `json:"servicetype"`
	SlotID      int    `json:"slotid"`
}

type RescheduleServiceRequest struct {
	SlotID int `json:"slotid"`
}

type RequestProvider struct{
	AssignedTo string `json:"assignedto"`
}