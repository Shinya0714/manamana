package rooting

import (
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
)

// type Owner struct {
// 	id         uint   `gorm:"primary_key"`
// 	createTime string `column:"create_time"`
// 	updateTime string `column:"update_time"`
// 	balance    int    `column:"balance"`
// }

type Data struct {
	List map[string]Item `json:"	"`
}

type Item struct {
	TargetCdString                           string `json:targetCdString`
	TargetPriceString                        string `json:targetPriceString`
	CompanyNameString                        string `json:companyNameString`
	BookBuildingString                       string `json:bookBuildingString`
	BookBuildingPossibleBoolStringForSbi     string `json:bookBuildingPossibleBoolStringForSbi`
	BookBuildingPossibleBoolStringForMizuho  string `json:bookBuildingPossibleBoolStringForMizuho`
	BookBuildingPossibleBoolStringForRakuten string `json:bookBuildingPossibleBoolStringForRakuten`
}

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

func sbiBookBuildingMap() map[string]string {

	driver, page := templateWebDriver()

	err := page.Navigate("https://www.sbisec.co.jp/ETGate")
	if err != nil {

		fmt.Println(err)
	}

	element := page.FindByXPath("//*[@id='user_input']/input")
	err = element.Fill(os.Getenv("SBI_USERNAME"))
	if err != nil {

		fmt.Println(err)
	}

	element = page.FindByXPath("//*[@id='password_input']/input")
	err = element.Fill(os.Getenv("SBI_LOGIN_PASSWORD"))
	if err != nil {

		fmt.Println(err)
	}

	page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()

	page.Navigate("https://site2.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETmgR001Control&_DataStoreID=DSWPLETmgR001Control&burl=search_domestic&dir=ipo%2F&file=stock_info_ipo.html&cat1=domestic&cat2=ipo&getFlg=on&int_pr1=150313_cmn_gnavi:6_dmenu_04")

	page.FindByXPath("/html/body/div[4]/div/table/tbody/tr/td[1]/div/div[10]/div/div/a/img").Click()

	time.Sleep(3 * time.Second)

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		targetCd := ""

		bookBuildingPossibleString := "false"

		target, err := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/div[2]/table[" + strconv.Itoa(i) + "]/tbody/tr/td/table/tbody/tr[2]/td[5]").Text()
		if err != nil {

			// NOOP
		} else {

			targetTitle, err := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/div[2]/table[" + strconv.Itoa(i) + "]/tbody/tr/td/table/tbody/tr[1]/td/table/tbody/tr/td[1]").Text()
			if err != nil {

				// NOOP
			} else {

				targetTitle = strings.ReplaceAll(targetTitle, "（株）", "")

				re := regexp.MustCompile(`（.+?）`)

				res := re.FindAllStringSubmatch(targetTitle, -1)[0]

				for _, v := range res {

					targetCd = strings.ReplaceAll(strings.ReplaceAll(v, "（", ""), "）", "")
				}
			}

			if strings.EqualFold(target, "") {

				bookBuildingPossibleString = "true"
			} else if strings.EqualFold(target, "取消   訂正") {

				bookBuildingPossibleString = "kanryo"
			} else {

				bookBuildingPossibleString = "false"
			}
		}

		m[targetCd] = bookBuildingPossibleString
	}

	err = driver.Stop()
	if err != nil {

		log.Println("sbiBookBuildingMap driver.Stop()", err)
	}

	return m
}

func mizuhoBookBuildingMap() map[string]string {

	driver, page := templateWebDriver()
	defer driver.Stop()

	page.Navigate("https://www.mizuho-sc.com/index.html")

	page.FindByXPath("/html/body/div[3]/div/header/div/div[5]/div/ul/li[7]/a/span").Click()

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_ID"))

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/div/ul/li[1]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_PASSWORD"))

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[2]/div/button").Click()

	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[1]/input[1]").Fill(os.Getenv("MIZUHO_1"))
	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[2]/input[1]").Fill(os.Getenv("MIZUHO_2"))

	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_3"))
	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[3]/label/input[1]").Fill(os.Getenv("MIZUHO_4"))
	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[5]/label/input[1]").Fill(os.Getenv("MIZUHO_5"))

	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[2]/div/button").Click()

	page.Navigate("https://mnc.mizuho-sc.com/web/rmfTrdStkIpoLstAction.do#IPOList")

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		bookBuildingPossibleString := "false"

		target, _ := page.AllByXPath("/html/body/div[1]/div[2]/main/div/div[4]/div[" + strconv.Itoa(i) + "]/div[1]/div/h3/span[1]/span").Text()

		bookBuildingStatus, _ := page.AllByXPath("/html/body/div[1]/div[2]/main/div/div[4]/div[" + strconv.Itoa(i) + "]/div[1]/div/div/a").Text()

		if strings.Contains(bookBuildingStatus, "抽選申込へ") {

			bookBuildingPossibleString = "true"
		} else if strings.Contains(bookBuildingStatus, "抽選申込取消へ") {

			bookBuildingPossibleString = "kanryo"
		} else {

			bookBuildingPossibleString = "false"
		}

		m[target] = bookBuildingPossibleString
	}

	return m
}

func rakutenBookBuildingMap() map[string]string {

	driver, page := templateWebDriver()
	defer driver.Stop()

	page.Navigate("https://www.rakuten-sec.co.jp/")

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/input[1]").Fill(os.Getenv("RAKUTEN_ID"))

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/div/input").Fill(os.Getenv("RAKUTEN_PASSWORD"))

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/button").Click()

	page.FindByXPath("/html/body/div[1]/div/div[5]/div/ul/li[3]/a/span").Click()

	page.FindByXPath("/html/body/div[1]/div/div[7]/div/ul/li[3]/a").Click()

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

	return m
}

func getSbiBalance(result chan string) {

	log.Printf("getSbiBalance start")
	defer log.Printf("getSbiBalance end")

	driver, page := templateWebDriver()
	defer driver.Stop()

	page.Navigate("https://www.sbisec.co.jp/ETGate")

	page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))

	page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))

	page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()

	time.Sleep(3 * time.Second)

	page.FindByXPath("/html/body/div[1]/div[1]/div[2]/div/ul/li[3]/a/img").Click()

	time.Sleep(3 * time.Second)

	kaitsukeKano2daysAfter, err := page.FindByXPath("/html/body/div[1]/table/tbody/tr/td[1]/table/tbody/tr[2]/td/table[1]/tbody/tr/td/form/table[2]/tbody/tr[1]/td[2]/table[4]/tbody/tr/td[1]/table[2]/tbody/tr[3]/td[2]/div").Text()
	if err != nil {

		log.Printf("getSbiBalance err: %v\n", err)

		kaitsukeKano2daysAfter = "読み込み失敗"
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
	defer driver.Stop()

	page.Navigate("https://www.mizuho-sc.com/index.html")

	page.FindByXPath("/html/body/div[3]/div/header/div/div[5]/div/ul/li[7]/a/span").Click()

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_ID"))

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/div/ul/li[1]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_PASSWORD"))

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[2]/div/button").Click()

	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[1]/input[1]").Fill(os.Getenv("MIZUHO_1"))
	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li[1]/div/label[2]/input[1]").Fill(os.Getenv("MIZUHO_2"))

	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_3"))
	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[3]/label/input[1]").Fill(os.Getenv("MIZUHO_4"))
	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/ul/li[2]/div[5]/label/input[1]").Fill(os.Getenv("MIZUHO_5"))

	page.FindByXPath("/html/body/div[1]/div[1]/main/div/div[2]/form/div[2]/div/button").Click()

	zandaka, err := page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/div/div[1]/div[2]/div[1]/div/div[3]/table[2]/tbody/tr/td/span").Text()
	if err != nil {

		fmt.Printf("getMizuhoBalance err: %v\n", err)

		zandaka = "読み込み失敗"
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
	defer driver.Stop()

	page.Navigate("https://www.rakuten-sec.co.jp/")

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/input[1]").Fill(os.Getenv("RAKUTEN_ID"))

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/div/input").Fill(os.Getenv("RAKUTEN_PASSWORD"))

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/button").Click()

	time.Sleep(5 * time.Second)

	zandaka, err := page.FindByXPath("/html/body/div[1]/div[2]/main/form[2]/div[2]/div[1]/div[2]/div/p[1]/span[1]").Text()
	if err != nil {

		fmt.Printf("getRakutenBalance err: %v\n", err)

		zandaka = "読み込み失敗"
	}

	result <- strings.Replace(string(zandaka), " 円", "", -1)
}

func GetBalance(c echo.Context) (err error) {

	channel := make(chan string)

	go getSbiBalance(channel)
	channelResult := <-channel

	go getMizuhoBalance(channel)
	channelResult2 := <-channel

	go getSmbcBalance(channel)
	channelResult3 := <-channel

	go getRakutenBalance(channel)
	channelResult4 := <-channel

	jsonMap := map[string]string{
		"sbiBalance":     channelResult,
		"mizuhoBalance":  channelResult2,
		"smbcBalance":    channelResult3,
		"rakutenBalance": channelResult4,
	}

	for k, v := range jsonMap {
		log.Printf("key: %s, value: %s\n", k, v)
	}

	return c.JSON(http.StatusOK, jsonMap)
}

func GetSchedule(c echo.Context) (err error) {

	log.Printf("getSchedule start")
	defer log.Printf("getSchedule end")

	driver, page := templateWebDriver()
	defer driver.Stop()

	sbiBookBuildingMap := sbiBookBuildingMap()
	mizuhoBookBuildingMap := mizuhoBookBuildingMap()
	rakutenBookBuildingMap := rakutenBookBuildingMap()

	page.Navigate("https://www.nikkei.com/markets/kigyo/ipo/money-schedule/")

	xpathStringForBookBuildingSpan := ""
	xpathStringForCompanyNameSpan := ""
	xpathStringForTargetCd := ""
	xpathStringForTargetPrice := ""

	bookBuildingString := ""
	companyNameString := ""
	targetCdString := ""
	targetPriceString := ""

	var companyNameStringList []string
	var targetCdStringList []string
	var targetPriceStringList []string
	var bookBuildingStringList []string
	var bookBuildingPossibleBoolListForSbi []string
	var bookBuildingPossibleBoolListForMizuho []string
	var bookBuildingPossibleBoolListForRakuten []string

	var data = Data{}
	data.List = map[string]Item{}

	items := []Item{}

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

			companyNameStringList = append(companyNameStringList, companyNameString)
			targetCdStringList = append(targetCdStringList, targetCdString)
			targetPriceStringList = append(targetPriceStringList, targetPriceString)
			bookBuildingStringList = append(bookBuildingStringList, bookBuildingString)

			if strings.EqualFold(sbiBookBuildingMap[targetCdString], "kanryo") {

				bookBuildingPossibleBoolStringForSbi = "kanryo"
			}

			if strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "kikanGai") {

				bookBuildingPossibleBoolStringForSbi = "kikanGai"
			}

			if strings.EqualFold(sbiBookBuildingMap[targetCdString], "true") && strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "kikanNai") {

				bookBuildingPossibleBoolStringForSbi = "true"
			}

			bookBuildingPossibleBoolListForSbi = append(bookBuildingPossibleBoolListForSbi, bookBuildingPossibleBoolStringForSbi)

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

			bookBuildingPossibleBoolListForMizuho = append(bookBuildingPossibleBoolListForMizuho, bookBuildingPossibleBoolStringForMizuho)

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

			bookBuildingPossibleBoolListForRakuten = append(bookBuildingPossibleBoolListForRakuten, bookBuildingPossibleBoolStringForRakuten)

			item := Item{TargetCdString: targetCdString, TargetPriceString: targetPriceString, CompanyNameString: companyNameString, BookBuildingString: bookBuildingString, BookBuildingPossibleBoolStringForSbi: bookBuildingPossibleBoolStringForSbi, BookBuildingPossibleBoolStringForMizuho: bookBuildingPossibleBoolStringForMizuho, BookBuildingPossibleBoolStringForRakuten: bookBuildingPossibleBoolStringForRakuten}
			items = append(items, item)
		}
	}

	// jsonエンコード
	outputJson, err := json.Marshal(&items)
	if err != nil {

		log.Printf("getSchedule err: %v\n", err)
	}

	jsonMap := map[string]string{"outputJson": string(outputJson)}

	return c.JSON(http.StatusOK, jsonMap)
}

func MizuhoBookBuilding(c echo.Context) (err error) {

	log.Printf("mizuhoBookBuilding start")
	defer log.Printf("mizuhoBookBuilding end")

	driver, page := templateWebDriver()
	defer driver.Stop()

	page.Navigate("https://www.mizuho-sc.com/index.html")

	page.FindByXPath("/html/body/div[3]/div/header/div/div[5]/div/ul/li[7]/a").Click()

	time.Sleep(3 * time.Second)

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[1]/div/ul/li/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_ID"))

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[1]/fieldset/div/div[2]/div/div/ul/li[1]/div[1]/label/input[1]").Fill(os.Getenv("MIZUHO_PASSWORD"))

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[2]/form/div[2]/div/button").Click()

	page.Navigate("https://mnc.mizuho-sc.com/web/rmfTrdStkIpoLstAction.do#IPOList")

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		target, _ := page.AllByXPath("/html/body/div[1]/div[2]/main/div/div[4]/div[" + strconv.Itoa(i) + "]/div[1]/div/h3/span[1]/span").Text()

		targetXpathString := "/html/body/div[1]/div[2]/main/div/div[4]/div[" + strconv.Itoa(i) + "]/div[1]/div/div/a"

		m[target] = targetXpathString
	}

	tickerSymbol := c.Param("tickerSymbol")

	fmt.Println(tickerSymbol)

	page.FindByXPath(m[tickerSymbol]).Click()

	time.Sleep(3 * time.Second)

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/table/tbody/tr/td/p/a").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/div[2]/div/button").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/table[1]/tbody/tr/td/p/a").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/table[2]/tbody/tr/td/p/a").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/table[3]/tbody/tr/td/p/a").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[3]/form/div[3]/div/button").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[4]/form/div[2]/fieldset/div/div[2]/div/ul/li[4]/div/input[2]").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[4]/form/div[3]/div/div/fieldset/input[2]").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[4]/form/div[4]/div[1]/button").Click()

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[5]/form/div[2]/div/div/div/fieldset/dl/dd/div[1]/div[1]/input[2]").Fill(os.Getenv("MIZUHO_TORIHIKI_PASSWORD"))

	page.FindByXPath("/html/body/div[1]/div[2]/main/div/div[5]/form/div[3]/div[1]/button").Click()

	resultString, _ := page.Title()

	return c.JSON(http.StatusOK, resultString)
}

func SbiBookBuilding(c echo.Context) (err error) {

	log.Printf("sbiBookBuilding start")
	defer log.Printf("sbiBookBuilding end")

	driver, page := templateWebDriver()
	defer driver.Stop()

	page.Navigate("https://www.sbisec.co.jp/ETGate")

	page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))

	page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))

	page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()

	page.Navigate("https://site2.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETmgR001Control&_DataStoreID=DSWPLETmgR001Control&burl=search_domestic&dir=ipo%2F&file=stock_info_ipo.html&cat1=domestic&cat2=ipo&getFlg=on&int_pr1=150313_cmn_gnavi:6_dmenu_04")

	page.FindByXPath("/html/body/div[4]/div/table/tbody/tr/td[1]/div/div[10]/div/div/a/img").Click()

	time.Sleep(3 * time.Second)

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		targetCd := ""

		targetXpathTableString := "false"

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

	page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/div[2]/table[" + m[tickerSymbol] + "]/tbody/tr/td/table/tbody/tr[2]/td[5]/a/img").Click()

	page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[6]/tbody/tr/td/table/tbody/tr/td[1]/table/tbody/tr[2]/td/input").Fill("100")

	page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[6]/tbody/tr/td/table/tbody/tr/td[2]/table/tbody/tr[2]/td[1]/input").Click()

	page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[8]/tbody/tr/td[1]/table/tbody/tr/td[2]/input").Fill("NFXQCBUM")

	page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[8]/tbody/tr/td[1]/table/tbody/tr/td[3]/input").Click()

	page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/form/table[7]/tbody/tr[2]/td/input[1]").Click()

	resultString, _ := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/table[3]/tbody/tr[1]/td/b").Text()

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
	defer driver.Stop()

	tickerSymbol := c.Param("tickerSymbol")

	fmt.Println(tickerSymbol)

	// 対象サイトに移動
	page.Navigate("https://www.rakuten-sec.co.jp/")

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/input[1]").Fill(os.Getenv("RAKUTEN_ID"))

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/div/input").Fill(os.Getenv("RAKUTEN_PASSWORD"))

	page.FindByXPath("/html/body/div[2]/section[1]/div/div[1]/div[3]/div[1]/form/button").Click()

	page.FindByXPath("/html/body/div[1]/div/div[5]/div/ul/li[3]/a/span").Click()

	page.FindByXPath("/html/body/div[1]/div/div[7]/div/ul/li[3]/a").Click()

	m := make(map[string]string)

	for i := 2; i < 50; i++ {

		targetCd, _ := page.AllByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table/tbody/tr[" + strconv.Itoa(i) + "]/td[1]/div/nobr").Text()

		i++

		targetXpath := "/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table/tbody/tr[" + strconv.Itoa(i) + "]/td[1]/div/nobr/a"

		m[targetCd] = targetXpath
	}

	page.FindByXPath(m[tickerSymbol]).Click()

	page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td[1]/div/table/tbody/tr/td/form/table[2]/tbody/tr/td/div/input[1]").Click()

	page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[2]/tbody/tr[1]/td[1]/div/nobr/input").Fill("100")

	page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[2]/tbody/tr[2]/td[1]/div/nobr/select").Select("成行")

	page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[3]/tbody/tr/td/div/input").Click()

	page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[2]/tbody/tr/td/table/tbody/tr/td[2]/input").Fill(os.Getenv("RAKUTEN_TORIHIKI_PASSWORD"))

	page.FindByXPath("/html/body/div[2]/div/div[1]/div/table/tbody/tr/td/div/table/tbody/tr/td/form/table[3]/tbody/tr/td/div/input[1]").Click()

	return c.JSON(http.StatusOK, "成功")
}

func templateWebDriver() (*agouti.WebDriver, *agouti.Page) {

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=1980,1200",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	driver.Start()

	page, err := driver.NewPage()
	if err != nil {

		panic(err)
	}

	return driver, page
}
