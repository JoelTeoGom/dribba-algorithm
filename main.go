package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Schedule struct {
	OpenFrom string `json:"openFrom"`
	OpenTo   string `json:"openTo"`
}

type Shop struct {
	ID        int        `json:"id"`
	Schedules []Schedule `json:"schedules"`
}

type Business struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Shops []Shop `json:"shops"`
}

func timeToMinutes(time string) int {
	var hour, minute int
	_, err := fmt.Sscanf(time, "%d:%d", &hour, &minute)
	if err != nil {
		log.Fatalf("Error parsejan temps: %v", err)
	}
	return hour*60 + minute
}

func minutesToTime(minutes int) string {
	hours := minutes / 60
	minutes = minutes % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func main() {
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error llegin el fitxer: %v", err)
	}

	var business Business

	err = json.Unmarshal(data, &business)
	if err != nil {
		log.Fatalf("Error parsejan el JSON: %v", err)
	}

	var timeSlots [1440]bool

	for _, shop := range business.Shops {
		for _, schedule := range shop.Schedules {
			open := timeToMinutes(schedule.OpenFrom)
			close := timeToMinutes(schedule.OpenTo)
			for i := open; i < close; i++ {
				timeSlots[i] = true
			}
		}
	}

	var result []Schedule

	hasStarted := false
	var newSchedule Schedule
	for i, timeSlot := range timeSlots {
		if !hasStarted && timeSlot {
			newSchedule.OpenFrom = minutesToTime(i)
			hasStarted = true
		} else if !timeSlot && hasStarted {
			newSchedule.OpenTo = minutesToTime(i)
			result = append(result, newSchedule)
			hasStarted = false
		}

	}
	if hasStarted {
		newSchedule.OpenTo = "23:59"
		result = append(result, newSchedule)
	}
	fmt.Println(result)
}


// Result:
// [
//   {
//     "openFrom": "09:20",
//     "openTo": "11:32"
//   },
//   {
//     "openFrom": "11:58",
//     "openTo": "21:30"
//   },
//   {
//     "openFrom": "21:31",
//     "openTo": "23:12"
//   }
// ]
