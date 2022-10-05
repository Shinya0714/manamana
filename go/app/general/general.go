package general

import (
	"fmt"
	"github.com/joho/godotenv"
	"strconv"
	"strings"
	"time"
)

// func GetStatusMap() map[string]int {

// 	m := map[string]int{
// 		"CANCEL":            0,
// 		"NO_HANDLING":       1,
// 		"AVAILABLE":         2,
// 		"WITHIN_THE_PERIOD": 3,
// 		"OUT_OF_TERM":       4,
// 		"RESERVED":          5,
// 		"NO_RESERVATION":    6,
// 	}

// 	return m
// }

func CheckBookoBuildingPossible(bookBuildingString string) string {

	bookoBuildingPossible := "false"

	t := time.Now()

	today := t.Format("20060102")

	if bookBuildingString != "---" {

		fromMonthint, _ := strconv.Atoi(strings.Split(strings.Split(bookBuildingString, "-")[0], "/")[0])
		fromDayint, _ := strconv.Atoi(strings.Split(strings.Split(bookBuildingString, "-")[0], "/")[1])

		toMonthint, _ := strconv.Atoi(strings.Split(strings.Split(bookBuildingString, "-")[1], "/")[0])
		toDayint, _ := strconv.Atoi(strings.Split(strings.Split(bookBuildingString, "-")[1], "/")[1])

		fromDate := strconv.Itoa(t.Year()) + fmt.Sprintf("%02d", fromMonthint) + fmt.Sprintf("%02d", fromDayint)
		toDate := strconv.Itoa(t.Year()) + fmt.Sprintf("%02d", toMonthint) + fmt.Sprintf("%02d", toDayint)

		toDayInt, _ := strconv.Atoi(today)
		fromDateInt, _ := strconv.Atoi(fromDate)
		toDateInt, _ := strconv.Atoi(toDate)

		if fromDateInt <= toDayInt && toDayInt <= toDateInt {

			bookoBuildingPossible = "kikanNai"
		} else {

			bookoBuildingPossible = "kikanGai"
		}
	}

	return bookoBuildingPossible
}

func LoadEnv() {

	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} else {

		fmt.Printf(".env読み込み完了")
	}
}
