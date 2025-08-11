package model

import "time"

type Status string

const(
	StatusPending Status = "pending"
	StatusApproved Status = "approved"
	StatusInProgress Status = "inProgress"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

type ServiceType string

const(
	Electrician ServiceType = "electrician"
	Plumber ServiceType = "plumber"
)

type ServiceRequest struct {
	RequestID   string
	ResidentID  string
	Status      Status
	TimeSlot    string
	StartTime   time.Time
	EndTime     time.Time
	ServiceType ServiceType
}
