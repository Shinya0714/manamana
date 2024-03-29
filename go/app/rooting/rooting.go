package rooting

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/lib/pq"

	"github.com/Shinya0714/manamana/go/app/general"

	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"github.com/sclevine/agouti"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// type Owner struct {
// 	id         uint   `gorm:"primary_key"`
// 	createTime string `column:"create_time"`
// 	updateTime string `column:"update_time"`
// 	balance    int    `column:"balance"`
// }

type Data struct {
	List map[string]Item `json:""`
}

//nolint:govet
type Item struct {
	TargetCdString                           string `json:targetCdString`
	TargetPriceString                        string `json:targetPriceString`
	CompanyNameString                        string `json:companyNameString`
	BookBuildingString                       string `json:bookBuildingString`
	BookBuildingPossibleBoolStringForSbi     string `json:bookBuildingPossibleBoolStringForSbi`
	BookBuildingPossibleBoolStringForMizuho  string `json:bookBuildingPossibleBoolStringForMizuho`
	BookBuildingPossibleBoolStringForRakuten string `json:bookBuildingPossibleBoolStringForRakuten`
}

const (
	trueString  = "true"
	falseString = "true"
)

// type Owner struct {
// 	id         uint   `gorm:"primary_key"`
// 	createTime string `column:"create_time"`
// 	updateTime string `column:"update_time"`
// 	balance    int    `column:"balance"`
// }

// func gormConnect() *gorm.DB {

// 	HOST := "db_container"
// 	PORT := "5432"
// 	USER := os.Getenv("POSTGRES_USER")
// 	PASSWORD := os.Getenv("POSTGRES_PASSWORD")
// 	DBNAME := os.Getenv("POSTGRES_DB")

// 	CONNECT := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", HOST, PORT, USER, DBNAME, PASSWORD)

// 	db, err := gorm.Open("postgres", CONNECT)
// 	if err != nil {

// 		panic(err.Error())
// 	}

// 	return db
// }

// db := gormConnect()
// defer db.Close()

// owner := []Owner{}

// // SELECT
// db.Find(&owner)

// for _, target := range owner {

// 	fmt.Println(target.id)
// 	fmt.Println(target.createTime)
// 	fmt.Println(target.updateTime)
// 	fmt.Println(target.balance)
// }

var count int

type Schedule struct {
	Title    string
	Location string
	Year     string
	Month    string
	Day      string
	Start    string
	End      string
}

func sbiBookBuildingMap() map[string]string {

	driver, page := templateWebDriver()

	err := page.Navigate("https://www.sbisec.co.jp/ETGate")
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()
	if err != nil {

		fmt.Println(err)
	}

	err = page.Navigate("https://site2.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETmgR001Control&_DataStoreID=DSWPLETmgR001Control&burl=search_domestic&dir=ipo%2F&file=stock_info_ipo.html&cat1=domestic&cat2=ipo&getFlg=on&int_pr1=150313_cmn_gnavi:6_dmenu_04")
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[4]/div/table/tbody/tr/td[1]/div/div[10]/div/div/a/img").Click()
	if err != nil {

		fmt.Println(err)
	}

	time.Sleep(3 * time.Second)

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		targetCd := ""

		bookBuildingPossibleString := "false"

		target, err := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/div[2]/table[" + strconv.Itoa(i) + "]/tbody/tr/td/table/tbody/tr[2]/td[5]").Text()
		if err != nil {

			// fmt.Println(err)
		} else {

			targetTitle, err := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/div[2]/table[" + strconv.Itoa(i) + "]/tbody/tr/td/table/tbody/tr[1]/td/table/tbody/tr/td[1]").Text()
			if err != nil {

				// fmt.Println(err)
			} else {

				targetTitle = strings.ReplaceAll(targetTitle, "（株）", "")

				re := regexp.MustCompile(`（.+?）`)

				res := re.FindAllStringSubmatch(targetTitle, -1)[0]

				for _, v := range res {

					targetCd = strings.ReplaceAll(strings.ReplaceAll(v, "（", ""), "）", "")
				}
			}

			if strings.EqualFold(target, "") {

				bookBuildingPossibleString = trueString
			} else if strings.EqualFold(target, "取消   訂正") {

				bookBuildingPossibleString = "kanryo"
			} else {

				bookBuildingPossibleString = falseString
			}
		}

		m[targetCd] = bookBuildingPossibleString
	}

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	return m
}

func mizuhoBookBuildingMap() map[string]string {

	driver, page := templateWebDriver()

	err := page.Navigate("https://www.mizuho-sc.com/index.html")
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[3]/div/header/div/div[5]/div/ul/li[7]/a/span").Click()
	if err != nil {

		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)

	err = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_ID"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/div/ul/li[1]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_PASSWORD"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[2]/div/button").Click()
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[1]/input[1]").Fill(os.Getenv("MIZUHO_1"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[2]/input[1]").Fill(os.Getenv("MIZUHO_2"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_3"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[3]/label/input[1]").Fill(os.Getenv("MIZUHO_4"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[5]/label/input[1]").Fill(os.Getenv("MIZUHO_5"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[2]/div/button").Click()
	if err != nil {

		fmt.Println(err)
	}

	err = page.Navigate("https://mnc.mizuho-sc.com/web/rmfTrdStkIpoLstAction.do#IPOList")
	if err != nil {

		fmt.Println(err)
	}

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		var bookBuildingPossibleString string

		target, _ := page.AllByXPath("/html/body/div[1]/div[2]/main/div/div[4]/div[" + strconv.Itoa(i) + "]/div[1]/div/h3/span[1]/span").Text()

		bookBuildingStatus, _ := page.AllByXPath("/html/body/div[1]/div[2]/main/div/div[4]/div[" + strconv.Itoa(i) + "]/div[1]/div/div/a").Text()

		if strings.Contains(bookBuildingStatus, "抽選申込へ") {

			bookBuildingPossibleString = "true"
		} else if strings.Contains(bookBuildingStatus, "抽選申込取消へ") {

			bookBuildingPossibleString = "kanryo"
		} else {

			bookBuildingPossibleString = "false"
		}

		fmt.Println(target)
		fmt.Println(bookBuildingPossibleString)

		m[target] = bookBuildingPossibleString
	}

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	return m
}

func rakutenBookBuildingMap() map[string]string {

	driver, page := templateWebDriver()

	err := page.Navigate("https://www.rakuten-sec.co.jp/")
	if err != nil {

		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)

	err = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/input[1]").Fill(os.Getenv("RAKUTEN_ID"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/div/input").Fill(os.Getenv("RAKUTEN_PASSWORD"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/button").Click()
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div/div[5]/div/ul/li[3]/a/span").Click()
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div/div[7]/div/ul/li[3]/a").Click()
	if err != nil {

		fmt.Println(err)
	}

	m := make(map[string]string)

	bookBuildingPossibleString := ""

	for i := 2; i < 50; i++ {

		targetCd, _ := page.AllByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table/tbody/tr[" + strconv.Itoa(i) + "]/td[1]/div/nobr").Text()

		i++

		targetStatus, _ := page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table/tbody/tr[" + strconv.Itoa(i) + "]/td[1]/div/nobr/a[1]").Text()

		if strings.Contains(targetStatus, "参加") {

			bookBuildingPossibleString = "true"
		} else if strings.Contains(targetStatus, "確認") {

			bookBuildingPossibleString = "kanryo"
		} else {

			bookBuildingPossibleString = "false"
		}

		m[targetCd] = bookBuildingPossibleString
	}

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	return m
}

func getSbiBalance(result chan string) {

	log.Printf("getSbiBalance start")
	defer log.Printf("getSbiBalance end")

	driver, page := templateWebDriver()

	err := page.Navigate("https://www.sbisec.co.jp/ETGate")
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()
	if err != nil {

		fmt.Println(err)
	}

	time.Sleep(3 * time.Second)

	err = page.FindByXPath("/html/body/div[1]/div[1]/div[2]/div/ul/li[3]/a/img").Click()
	if err != nil {

		fmt.Println(err)
	}

	time.Sleep(3 * time.Second)

	kaitsukeKano2daysAfter, err := page.FindByXPath("/html/body/div[1]/table/tbody/tr/td[1]/table/tbody/tr[2]/td/table[1]/tbody/tr/td/form/table[2]/tbody/tr[1]/td[2]/table[4]/tbody/tr/td[1]/table[2]/tbody/tr[3]/td[2]/div").Text()
	if err != nil {

		log.Printf("getSbiBalance err: %v\n", err)

		kaitsukeKano2daysAfter = "読み込み失敗"
	}

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	result <- kaitsukeKano2daysAfter
}

// func getDaiwaBalance(c echo.Context) (err error) {

// 	driver := agouti.ChromeDriver(

// 		agouti.ChromeOptions("args", []string{

// 			"--headless",
// 			"--window-size=300,1200",
// 			"--blink-settings=imagesEnabled=false",
// 			"--disable-gpu",
// 			"no-sandbox",
// 		}),
// 	)

// 	defer driver.Stop()
// 	driver.Start()

// 	page, err := driver.NewPage()
// 	if err != nil {

// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		return
// 	}

// 	// 対象サイトに移動
// 	page.Navigate("https://www.daiwa.co.jp/PCC/HomeTrade/Account/m8301.html")

// 	time.Sleep(3 * time.Second)

// 	page.FindByName("@PM-1@").Fill(os.Getenv("DAIWA_SHITENCD"))

// 	page.FindByName("@PM-2@").Fill(os.Getenv("DAIWA_KOZANUMBER"))

// 	page.FindByName("@PM-3@").Fill(os.Getenv("DAIWA_PASSWORD"))

// 	page.FindByXPath("//*[@id='CONTENT']/div[1]/div[2]/form/div[2]/input").Click()

// 	page.FindByXPath("//*[@id='menuTabsetHead']/form/table/tbody/tr/td/table/tbody/tr/td[7]/div[2]/a").Click()

// 	zandaka, err := page.FindByXPath("//*[@id='cTable']/tbody/tr[1]/td/table[4]/tbody/tr[2]/td[1]").Text()
// 	if err != nil {

// 		fmt.Printf("err: %v\n", err)
// 	}

// 	c.JSON(http.StatusOK, "買付余力"+zandaka)

// 	return
// }

func getMizuhoBalance(result chan string) {

	log.Printf("getMizuhoBalance start")
	defer log.Printf("getMizuhoBalance end")

	driver, page := templateWebDriver()

	err := page.Navigate("https://www.mizuho-sc.com/index.html")
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[3]/div/header/div/div[5]/div/ul/li[7]/a/span").Click()
	if err != nil {

		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)

	err = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_ID"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/div/ul/li[1]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_PASSWORD"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[2]/div/button").Click()
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[1]/input[1]").Fill(os.Getenv("MIZUHO_1"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[2]/input[1]").Fill(os.Getenv("MIZUHO_2"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_3"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[3]/label/input[1]").Fill(os.Getenv("MIZUHO_4"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[5]/label/input[1]").Fill(os.Getenv("MIZUHO_5"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[2]/div/button").Click()
	if err != nil {

		fmt.Println(err)
	}

	zandaka, err := page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/div/div[1]/div[2]/div[1]/div/div[3]/table[2]/tbody/tr/td/span").Text()
	if err != nil {

		fmt.Printf("getMizuhoBalance err: %v\n", err)

		zandaka = "読み込み失敗"
	}

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	result <- zandaka
}

func getSmbcBalance(result chan string) {

	log.Printf("getSmbcBalance start")
	defer log.Printf("getSmbcBalance end")

	zandaka, err := exec.Command("python3", "/Users/Owner/manamana/go/app/rooting/smbc_zandaka.py").CombinedOutput()
	if err != nil {

		fmt.Println(err)
	}

	fmt.Println(strings.Replace(string(zandaka), "万円", "0,000", -1))

	result <- strings.Replace(string(zandaka), "万円", "0,000", -1)
}

func getRakutenBalance(result chan string) {

	log.Printf("getRakutenBalance start")
	defer log.Printf("getRakutenBalance end")

	driver, page := templateWebDriver()

	err := page.Navigate("https://www.rakuten-sec.co.jp/")
	if err != nil {

		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)

	err = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/input[1]").Fill(os.Getenv("RAKUTEN_ID"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/div/input").Fill(os.Getenv("RAKUTEN_PASSWORD"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/button").Click()
	if err != nil {

		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)

	zandaka, err := page.FindByXPath("/html/body/div[1]/div[2]/main/form[2]/div[2]/div[1]/div[2]/div/p[1]/span[1]").Text()
	if err != nil {

		fmt.Printf("getRakutenBalance err: %v\n", err)

		zandaka = "読み込み失敗"
	}

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	result <- strings.Replace(string(zandaka), " 円", "", -1)
}

func GetBalance(c echo.Context) (err error) {

	count = 0

	count = 1 * 10
	time.Sleep(time.Second)

	channel := make(chan string)

	count = 2 * 10
	time.Sleep(time.Second)

	go getSbiBalance(channel)
	channelResult := <-channel

	count = 3 * 10
	time.Sleep(time.Second)

	go getMizuhoBalance(channel)
	channelResult2 := <-channel

	count = 4 * 10
	time.Sleep(time.Second)

	go getSmbcBalance(channel)
	channelResult3 := <-channel

	count = 5 * 10
	time.Sleep(time.Second)

	go getRakutenBalance(channel)
	channelResult4 := <-channel

	count = 6 * 10
	time.Sleep(time.Second)

	jsonMap := map[string]string{
		"sbiBalance":     channelResult,
		"mizuhoBalance":  channelResult2,
		"smbcBalance":    channelResult3,
		"rakutenBalance": channelResult4,
	}

	count = 7 * 10
	time.Sleep(time.Second)

	for k, v := range jsonMap {
		log.Printf("key: %s, value: %s\n", k, v)
	}

	count = 9 * 10
	time.Sleep(time.Second)

	count = 10 * 10
	time.Sleep(time.Second)

	return c.JSON(http.StatusOK, jsonMap)
}

func GetSchedule(c echo.Context) (err error) {

	count = 0

	log.Printf("getSchedule start")
	defer log.Printf("getSchedule end")

	count = 1 * 10
	time.Sleep(time.Second)

	driver, page := templateWebDriver()

	count = 2 * 10
	time.Sleep(time.Second)

	sbiBookBuildingMap := sbiBookBuildingMap()
	mizuhoBookBuildingMap := mizuhoBookBuildingMap()
	rakutenBookBuildingMap := rakutenBookBuildingMap()

	count = 3 * 10
	time.Sleep(time.Second)

	error := page.Navigate("https://www.nikkei.com/markets/kigyo/ipo/money-schedule/")
	if error != nil {

		fmt.Println(error)
	}

	count = 4 * 10
	time.Sleep(time.Second)

	xpathStringForBookBuildingSpan := ""
	xpathStringForCompanyNameSpan := ""
	xpathStringForTargetCd := ""
	xpathStringForTargetPrice := ""

	bookBuildingString := ""
	companyNameString := ""
	targetCdString := ""
	targetPriceString := ""

	count = 5 * 10
	time.Sleep(time.Second)

	var data = Data{}
	data.List = map[string]Item{}

	items := []Item{}

	count = 6 * 10
	time.Sleep(time.Second)

	for i := 1; i <= 50; i++ {

		bookBuildingPossibleBoolStringForSbi := "false"
		bookBuildingPossibleBoolStringForMizuho := "false"
		bookBuildingPossibleBoolStringForRakuten := "false"

		xpathStringForCompanyNameSpan = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[2]", i)
		xpathStringForBookBuildingSpan = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[3]", i)
		xpathStringForTargetCd = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[1]/a", i)
		xpathStringForTargetPrice = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[4]", i)

		companyNameString, _ = page.FindByXPath(xpathStringForCompanyNameSpan).Text()
		targetCdString, _ = page.FindByXPath(xpathStringForTargetCd).Text()
		targetPriceString, _ = page.FindByXPath(xpathStringForTargetPrice).Text()
		bookBuildingString, err = page.FindByXPath(xpathStringForBookBuildingSpan).Text()
		if err != nil {

			fmt.Printf("getSchedule err: %v\n", err)

			break
		} else {

			if strings.EqualFold(sbiBookBuildingMap[targetCdString], "kanryo") {

				bookBuildingPossibleBoolStringForSbi = "kanryo"
			}

			if strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "kikanGai") {

				bookBuildingPossibleBoolStringForSbi = "kikanGai"
			}

			if strings.EqualFold(sbiBookBuildingMap[targetCdString], "true") && strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "kikanNai") {

				bookBuildingPossibleBoolStringForSbi = "true"
			}

			// みずほ
			if strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "kikanGai") {

				bookBuildingPossibleBoolStringForMizuho = "kikanGai"
			}

			if strings.EqualFold(mizuhoBookBuildingMap[targetCdString], "kanryo") {

				bookBuildingPossibleBoolStringForMizuho = "kanryo"
			}

			if strings.EqualFold(mizuhoBookBuildingMap[targetCdString], "true") && strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "kikanNai") {

				bookBuildingPossibleBoolStringForMizuho = "true"
			}

			// 楽天
			if strings.EqualFold(rakutenBookBuildingMap[targetCdString], "kanryo") {

				bookBuildingPossibleBoolStringForRakuten = "kanryo"
			}

			if strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "kikanGai") {

				bookBuildingPossibleBoolStringForRakuten = "kikanGai"
			}

			if strings.EqualFold(rakutenBookBuildingMap[targetCdString], "true") && strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "kikanNai") {

				bookBuildingPossibleBoolStringForRakuten = "true"
			}

			item := Item{TargetCdString: targetCdString, TargetPriceString: targetPriceString, CompanyNameString: companyNameString, BookBuildingString: bookBuildingString, BookBuildingPossibleBoolStringForSbi: bookBuildingPossibleBoolStringForSbi, BookBuildingPossibleBoolStringForMizuho: bookBuildingPossibleBoolStringForMizuho, BookBuildingPossibleBoolStringForRakuten: bookBuildingPossibleBoolStringForRakuten}
			items = append(items, item)
		}
	}

	count = 7 * 10
	time.Sleep(time.Second)

	// jsonエンコード
	outputJson, err := json.Marshal(&items)
	if err != nil {

		log.Printf("getSchedule err: %v\n", err)
	}

	count = 8 * 10
	time.Sleep(time.Second)

	jsonMap := map[string]string{"outputJson": string(outputJson)}

	count = 9 * 10
	time.Sleep(time.Second)

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	count = 10 * 10
	time.Sleep(time.Second)

	return c.JSON(http.StatusOK, jsonMap)
}

func MizuhoBookBuilding(c echo.Context) (err error) {

	log.Printf("mizuhoBookBuilding start")
	defer log.Printf("mizuhoBookBuilding end")

	driver, page := templateWebDriver()

	error := page.Navigate("https://www.mizuho-sc.com/index.html")
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[3]/div/header/div/div[5]/div/ul/li[7]/a").Click()
	if error != nil {

		fmt.Println(error)
	}

	time.Sleep(3 * time.Second)

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_ID"))
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/div/ul/li[1]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_PASSWORD"))
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[2]/div/button").Click()
	if error != nil {

		fmt.Println(error)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[1]/input[1]").Fill(os.Getenv("MIZUHO_1"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[2]/input[1]").Fill(os.Getenv("MIZUHO_2"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_3"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[3]/label/input[1]").Fill(os.Getenv("MIZUHO_4"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[5]/label/input[1]").Fill(os.Getenv("MIZUHO_5"))
	if err != nil {

		fmt.Println(err)
	}

	err = page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[2]/div/button").Click()
	if err != nil {

		fmt.Println(err)
	}

	error = page.Navigate("https://mnc.mizuho-sc.com/web/rmfTrdStkIpoLstAction.do#IPOList")
	if error != nil {

		fmt.Println(error)
	}

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		target, _ := page.AllByXPath("/html/body/div[1]/div[2]/main/div/div[4]/div[" + strconv.Itoa(i) + "]/div[1]/div/h3/span[1]/span").Text()

		targetXpathString := "/html/body/div[1]/div[2]/main/div/div[4]/div[" + strconv.Itoa(i) + "]/div[1]/div/div/a"

		m[target] = targetXpathString
	}

	tickerSymbol := c.Param("tickerSymbol")

	fmt.Println(tickerSymbol)

	error = page.FindByXPath(m[tickerSymbol]).Click()
	if error != nil {

		fmt.Println(error)
	}

	time.Sleep(3 * time.Second)

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/table/tbody/tr/td/p/a").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/div[2]/div/button").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/table[1]/tbody/tr/td/p/a").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/table[2]/tbody/tr/td/p/a").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/table[3]/tbody/tr/td/p/a").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/div[3]/div/button").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[4]/form/div[2]/fieldset/div/div[2]/div/ul/li[4]/div/input[2]").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[4]/form/div[3]/div/div/fieldset/input[2]").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[4]/form/div[4]/div[1]/button").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[5]/form/div[2]/div/div/div/fieldset/dl/dd/div[1]/div[1]/input[2]").Fill(os.Getenv("MIZUHO_TORIHIKI_PASSWORD"))
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[5]/form/div[3]/div[1]/button").Click()
	if error != nil {
		for i, v := range os.Args {
			fmt.Printf("args[%d] -> %s\n", i, v)
		}
		fmt.Println(error)
	}

	resultString, _ := page.Title()

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	addGoogleCalendar()

	return c.JSON(http.StatusOK, resultString)
}

func SbiBookBuilding(c echo.Context) (err error) {

	log.Printf("sbiBookBuilding start")
	defer log.Printf("sbiBookBuilding end")

	driver, page := templateWebDriver()

	error := page.Navigate("https://www.sbisec.co.jp/ETGate")
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.Navigate("https://site2.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETmgR001Control&_DataStoreID=DSWPLETmgR001Control&burl=search_domestic&dir=ipo%2F&file=stock_info_ipo.html&cat1=domestic&cat2=ipo&getFlg=on&int_pr1=150313_cmn_gnavi:6_dmenu_04")
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[4]/div/table/tbody/tr/td[1]/div/div[10]/div/div/a/img").Click()
	if error != nil {

		fmt.Println(error)
	}

	time.Sleep(3 * time.Second)

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		targetCd := ""

		var targetXpathTableString string

		target, _ := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/div[2]/table[" + strconv.Itoa(i) + "]/tbody/tr/td/table/tbody/tr[2]/td[5]").Text()

		targetTitle, _ := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/div[2]/table[" + strconv.Itoa(i) + "]/tbody/tr/td/table/tbody/tr[1]/td/table/tbody/tr/td[1]").Text()

		targetTitle = strings.ReplaceAll(targetTitle, "（株）", "")

		re := regexp.MustCompile(`（.+?）`)

		res := re.FindAllStringSubmatch(targetTitle, -1)[0]

		for _, v := range res {

			targetCd = strings.ReplaceAll(strings.ReplaceAll(v, "（", ""), "）", "")
		}

		if target == "" {

			targetXpathTableString = strconv.Itoa(i)
		} else {

			targetXpathTableString = "false"
		}

		m[targetCd] = targetXpathTableString
	}

	tickerSymbol := c.Param("tickerSymbol")

	error = page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/div[2]/table[" + m[tickerSymbol] + "]/tbody/tr/td/table/tbody/tr[2]/td[5]/a/img").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[6]/tbody/tr/td/table/tbody/tr/td[1]/table/tbody/tr[2]/td/input").Fill("100")
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[6]/tbody/tr/td/table/tbody/tr/td[2]/table/tbody/tr[2]/td[1]/input").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[8]/tbody/tr/td[1]/table/tbody/tr/td[2]/input").Fill("NFXQCBUM")
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[8]/tbody/tr/td[1]/table/tbody/tr/td[3]/input").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/form/table[7]/tbody/tr[2]/td/input[1]").Click()
	if error != nil {

		fmt.Println(error)
	}

	resultString, _ := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/table[3]/tbody/tr[1]/td/b").Text()

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, resultString)
}

func SmbcBookBuilding(c echo.Context) (err error) {

	log.Printf("SmbcBookBuilding start")
	defer log.Printf("SmbcBookBuilding end")

	resultString, err := exec.Command("python3", "/Users/Owner/manamana/go/app/rooting/smbc.py").CombinedOutput()
	if err != nil {

		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, resultString)
}

func RakutenBookBuilding(c echo.Context) (err error) {

	log.Printf("RakutenBookBuilding start")
	defer log.Printf("RakutenBookBuilding end")

	driver, page := templateWebDriver()

	tickerSymbol := c.Param("tickerSymbol")

	fmt.Println(tickerSymbol)

	error := page.Navigate("https://www.rakuten-sec.co.jp/")
	if error != nil {

		fmt.Println(error)
	}

	time.Sleep(5 * time.Second)

	error = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/input[1]").Fill(os.Getenv("RAKUTEN_ID"))
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/div/input").Fill(os.Getenv("RAKUTEN_PASSWORD"))
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/button").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div/div[5]/div/ul/li[3]/a/span").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[1]/div/div[7]/div/ul/li[3]/a").Click()
	if error != nil {

		fmt.Println(error)
	}

	m := make(map[string]string)

	for i := 2; i < 50; i++ {

		targetCd, _ := page.AllByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table/tbody/tr[" + strconv.Itoa(i) + "]/td[1]/div/nobr").Text()

		i++

		targetXpath := "/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table/tbody/tr[" + strconv.Itoa(i) + "]/td[1]/div/nobr/a"

		m[targetCd] = targetXpath
	}

	error = page.FindByXPath(m[tickerSymbol]).Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td[1]/div/table/tbody/tr/td/form/table[2]/tbody/tr/td/div/input[1]").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[2]/tbody/tr[1]/td[1]/div/nobr/input").Fill("100")
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[2]/tbody/tr[2]/td[1]/div/nobr/select").Select("成行")
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[3]/tbody/tr/td/div/input").Click()
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[2]/tbody/tr/td/table/tbody/tr/td[2]/input").Fill(os.Getenv("RAKUTEN_TORIHIKI_PASSWORD"))
	if error != nil {

		fmt.Println(error)
	}

	error = page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[3]/tbody/tr/td/div/input[1]").Click()
	if error != nil {

		fmt.Println(error)
	}

	err = driver.Stop()
	if err != nil {

		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, "成功")
}

func templateWebDriver() (*agouti.WebDriver, *agouti.Page) {

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			// "--headless",
			"--window-size=1980,1200",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	error := driver.Start()
	if error != nil {

		fmt.Println(error)
	}

	page, err := driver.NewPage()
	if err != nil {

		panic(err)
	}

	return driver, page
}

func ProgressFunc(c echo.Context) (err error) {

	fmt.Println(count)

	return c.JSON(http.StatusOK, count)
}

func getClient(config *oauth2.Config) *http.Client {

	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func SetSchedule(scheduleMap map[string]string) *Schedule {
	return &Schedule{
		Title:    scheduleMap["Title"],
		Location: scheduleMap["Location"],
		Year:     scheduleMap["Year"],
		Month:    scheduleMap["Month"],
		Day:      scheduleMap["Day"],
		Start:    scheduleMap["Start"],
		End:      scheduleMap["End"],
	}
}

func (schedule Schedule) createEventData() *calendar.Event {

	start_datatime := schedule.Year + "-" + schedule.Month + "-" + schedule.Day + "T" + schedule.Start + ":00:00+09:00"
	end_datatime := schedule.Year + "-" + schedule.Month + "-" + schedule.Day + "T" + schedule.End + ":00:00+09:00"

	event := &calendar.Event{
		Summary:  schedule.Title,
		Location: schedule.Location,
		Start: &calendar.EventDateTime{
			DateTime: start_datatime,
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			DateTime: end_datatime,
			TimeZone: "Asia/Tokyo",
		},
	}

	return event
}

func readFileByCredentialsJson() []byte {

	b, err := os.ReadFile("/Users/Owner/manamana/go/app/rooting/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	return b
}

func addGoogleCalendar() {

	config, err := google.ConfigFromJSON(readFileByCredentialsJson(), calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	m := map[string]string{}

	m["Title"] = "testのてすと"
	m["Location"] = "locationのてすと"
	m["Year"] = "2022"
	m["Month"] = "12"
	m["Day"] = "25"
	m["Start"] = "13"
	m["End"] = "14"

	schedule := SetSchedule(m)

	calendarIdString := os.Getenv("TARGET_ID_FOR_GOOGLE_CALENDAR_API")

	_, err = srv.Events.Insert(calendarIdString, schedule.createEventData()).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}
