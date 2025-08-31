package utils

import (
	"os"
	"testing"

	"github.com/MananLed/majorProjectSMS/internal/model"
)

func withStdin(input string, fn func()) {
    r, w, _ := os.Pipe()


    oldStdin := os.Stdin
    defer func() {
        os.Stdin = oldStdin
    }()

    w.Write([]byte(input))
    w.Close()

    os.Stdin = r

    fn()
}

func TestFilterBookedSlots(t *testing.T) {
	slots := GenerateTimeSlots()
	reqSlot := slots[1] 

	requests := []model.ServiceRequest{
		{
			ServiceType: model.Electrician,
			StartTime:   reqSlot.StartTime,
			EndTime:     reqSlot.EndTime,
		},
	}

	filtered := FilterBookedSlots(slots, requests, model.Electrician)
	if len(filtered) != len(slots)-1 {
		t.Fatalf("expected %d slots after filtering, got %d", len(slots)-1, len(filtered))
	}

	filtered = FilterBookedSlots(slots, requests, model.Plumber)
	if len(filtered) != len(slots) {
		t.Fatalf("expected %d slots when booking different service type, got %d", len(slots), len(filtered))
	}
}

func TestSlotBooking_ValidInput(t *testing.T) {
    slots := GenerateTimeSlots()

    withStdin("1\n", func() {
        slot, err := SlotBooking(slots)
        if err != nil {
            t.Fatalf("expected no error, got %v", err)
        }
        if slot != slots[0] {
            t.Fatalf("expected slot %v, got %v", slots[0], slot)
        }
    })
}

func TestSlotBooking_InvalidChoice(t *testing.T) {
    slots := GenerateTimeSlots()

    cases := []string{
        "0\n",       
        "999\n",      
        "not-a-number", 
    }

    for _, input := range cases {
        withStdin(input, func() {
            _, err := SlotBooking(slots)
            if err == nil {
                t.Fatalf("expected error for input %q, got nil", input)
            }
        })
    }
}