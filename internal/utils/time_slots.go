package utils

import (
	"fmt"
	"time"
	"errors"
	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/fatih/color"
	"github.com/MananLed/majorProjectSMS/internal/model"
)

type TimeSlot struct {
	Label     string
	StartTime time.Time
	EndTime   time.Time
}

func GenerateTimeSlots() []TimeSlot {
	slots := []TimeSlot{}
	start := time.Date(0, 1, 1, constants.StartTimeOfService, 0, 0, 0, time.Local)
	end := time.Date(0, 1, 1, constants.EndTimeOfService, 0, 0, 0, time.Local)

	for current := start; current.Add(time.Minute*constants.TimeLimitOfSlot).Before(end) || current.Add(time.Minute*constants.TimeLimitOfSlot).Equal(end); current = current.Add(time.Minute * constants.TimeLimitOfSlot) {
		endTime := current.Add(time.Minute * constants.TimeLimitOfSlot)

		slots = append(slots, TimeSlot{
			StartTime: current,
			EndTime:   endTime,
			Label:     fmt.Sprintf("%s - %s", current.Format("3:04 PM"), endTime.Format("3:04 PM")),
		})
	}

	return slots
}

func FilterBookedSlots(slots []TimeSlot, requests []model.ServiceRequest, serviceType model.ServiceType) []TimeSlot {
	filtered := []TimeSlot{}

slotLoop:
	for _, slot := range slots {
		for _, req := range requests {
			if req.ServiceType == serviceType && slot.StartTime.Equal(req.StartTime) && slot.EndTime.Equal(req.EndTime) {
				continue slotLoop
			}
		}
		filtered = append(filtered, slot)
	}

	return filtered
}

func SlotBooking(slots []TimeSlot) (TimeSlot, error) {
	for i, s := range slots {
		fmt.Printf("%d. %s\n", i+1, s.Label)
	}

	var choice int
	fmt.Print(color.CyanString("Enter slot number: "))
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(slots) {
		return TimeSlot{}, errors.New("invalid slot choice")
	}

	return slots[choice-1], nil
}
