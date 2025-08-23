package model

import "fmt"

type Society struct {
	Floors int
}

func (s *Society) initFloors() {
	s.Floors = 8
}

func GetNumberOfFloors() {
	var soc Society

	soc.initFloors()
	fmt.Printf("Number of Floors: %d\n", soc.Floors)
}
