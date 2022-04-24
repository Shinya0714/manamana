package general

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

			bookoBuildingPossible = "true"
		}
	}

	fmt.Printf("generalを更新出来てるかテスト")

	return bookoBuildingPossible
}
