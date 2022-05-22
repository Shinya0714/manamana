package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	"github.com/Shinya0714/manamana/go/app/general"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sclevine/agouti"
	"golang.org/x/net/websocket"
)

type Owner struct {
	id         uint   `gorm:"primary_key"`
	createTime string `column:"create_time"`
	updateTime string `column:"update_time"`
	balance    int    `column:"balance"`
}

func Handler() {

	loadEnv()

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

	e := echo.New()
	e.Use(middleware.CORS())

	// ルーティング
	e.GET("/sbiBalance", getSbiBalance)
	e.GET("/mizuhoBalance", getMizuhoBalance)
	e.GET("/schedule", getSchedule)
	e.GET("/sbiBookBuilding/:tickerSymbol", sbiBookBuilding)
	e.GET("/mizuhoBookBuilding/:tickerSymbol", mizuhoBookBuilding)
	e.GET("/ws", msgHandler)

	// local サーバー
	e.Logger.Fatal(e.Start(":8000"))

	return
}

func loadEnv() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}

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

func sbiBookBuildingMap() map[string]string {

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=300,1200",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	defer driver.Stop()
	driver.Start()

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}

	// 対象サイトに移動
	page.Navigate("https://www.sbisec.co.jp/ETGate")

	// // ユーザーネーム
	// page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))

	// // パスワード
	// page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))

	// 「ログイン」
	page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()

	// 「ブックビルディング情報」
	page.Navigate("https://site2.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETmgR001Control&_DataStoreID=DSWPLETmgR001Control&burl=search_domestic&dir=ipo%2F&file=stock_info_ipo.html&cat1=domestic&cat2=ipo&getFlg=on&int_pr1=150313_cmn_gnavi:6_dmenu_04")

	// 「新規上場株式ブックビルディング／購入意思表示」
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

			if target == "" {

				bookBuildingPossibleString = "true"
			} else {

				bookBuildingPossibleString = "false"
			}
		}

		m[targetCd] = bookBuildingPossibleString
	}

	return m
}

func mizuhoBookBuildingMap() map[string]string {

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=1920,1080",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	defer driver.Stop()
	driver.Start()

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	// 対象サイトに移動
	page.Navigate("https://netclub.mizuho-sc.com/mnc/login?rt_bn=sc_top_hd_login")

	time.Sleep(3 * time.Second)

	page.FindByXPath("/html/body/header[1]/div/div[1]/div/div/div[2]/ul/li[2]").Click()

	page.FindByXPath("//*[@id='form01']/p/span/input").Click()

	page.Navigate("https://netclub.mizuho-sc.com/mnc/tr/ipopo?6")

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		bookBuildingPossibleString := "false"

		target, err := page.AllByXPath("/html/body/div[2]/div[4]/div[2]/span[2]/span/table/tbody/tr[" + strconv.Itoa(i) + "]/td[2]/span").Text()
		if err != nil {

			// NOOP
		}

		bookBuildingStatus, err := page.AllByXPath("/html/body/div[2]/div[4]/div[2]/span[2]/span/table/tbody/tr[" + strconv.Itoa(i) + "]/td[1]/ul/li").Text()
		if err != nil {

			// NOOP
		}

		if strings.EqualFold(bookBuildingStatus, "申込") {

			bookBuildingPossibleString = "true"
		}

		m[target] = bookBuildingPossibleString
	}

	return m
}

func getSbiBalance(c echo.Context) (err error) {

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=300,1200",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	defer driver.Stop()
	driver.Start()

	fmt.Println("driver読み込み完了")

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// 対象サイトに移動
	page.Navigate("https://www.sbisec.co.jp/ETGate")

	// // ユーザーネーム
	// page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))

	// // パスワード
	// page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))

	// 「ログイン」
	page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()

	time.Sleep(3 * time.Second)

	page.FindByXPath("/html/body/div[1]/div[1]/div[2]/div/ul/li[3]/a/img").Click()

	time.Sleep(3 * time.Second)

	kaitsukeKano2daysAfter, err := page.FindByXPath("/html/body/div[1]/table/tbody/tr/td[1]/table/tbody/tr[2]/td/table[1]/tbody/tr/td/form/table[2]/tbody/tr[1]/td[2]/table[4]/tbody/tr/td[1]/table[2]/tbody/tr[3]/td[2]/div").Text()
	if err != nil {

		fmt.Printf("err: %v\n", err)
	}

	c.JSON(http.StatusOK, kaitsukeKano2daysAfter)

	return
}

func getDaiwaBalance(c echo.Context) (err error) {

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=300,1200",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	defer driver.Stop()
	driver.Start()

	fmt.Println("driver読み込み完了")

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// 対象サイトに移動
	page.Navigate("https://www.daiwa.co.jp/PCC/HomeTrade/Account/m8301.html")

	time.Sleep(3 * time.Second)

	page.FindByName("@PM-1@").Fill(os.Getenv("DAIWA_SHITENCD"))

	page.FindByName("@PM-2@").Fill(os.Getenv("DAIWA_KOZANUMBER"))

	page.FindByName("@PM-3@").Fill(os.Getenv("DAIWA_PASSWORD"))

	page.FindByXPath("//*[@id='CONTENT']/div[1]/div[2]/form/div[2]/input").Click()

	page.FindByXPath("//*[@id='menuTabsetHead']/form/table/tbody/tr/td/table/tbody/tr/td[7]/div[2]/a").Click()

	zandaka, err := page.FindByXPath("//*[@id='cTable']/tbody/tr[1]/td/table[4]/tbody/tr[2]/td[1]").Text()
	if err != nil {

		fmt.Printf("err: %v\n", err)
	}

	c.JSON(http.StatusOK, "買付余力"+zandaka)

	return
}

func getMizuhoBalance(c echo.Context) (err error) {

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=300,1200",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	defer driver.Stop()
	driver.Start()

	fmt.Println("driver読み込み完了")

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// 対象サイトに移動
	page.Navigate("https://netclub.mizuho-sc.com/mnc/login?rt_bn=sc_top_hd_login")

	fmt.Println(page.URL())

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/header[1]/div/div[1]/div/div/div[2]/ul/li[2]").Click()

	// page.FindByXPath("//*[@id='IDInputKB']").Fill(os.Getenv("MIZUHO_ID"))

	// page.FindByXPath("//*[@id='PWInputKB']").Fill(os.Getenv("MIZUHO_PASSWORD"))

	page.FindByXPath("//*[@id='form01']/p/span/input").Click()

	time.Sleep(5 * time.Second)

	zandaka, err := page.FindByXPath("/html/body/div[2]/div[3]/div/div/div/div[1]/span[2]/div/table/tbody/tr/td[2]/table/tbody/tr/td/span/strong/span").Text()
	if err != nil {

		fmt.Printf("err: %v\n", err)
	}

	c.JSON(http.StatusOK, zandaka)

	return
}

func getSchedule(c echo.Context) (err error) {

	log.Println("getSchedule")

	sbiBookBuildingMap := sbiBookBuildingMap()
	mizuhoBookBuildingMap := mizuhoBookBuildingMap()

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=1920,1080",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	defer driver.Stop()
	driver.Start()

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	page.Navigate("https://www.nikkei.com/markets/kigyo/ipo/money-schedule/")

	xpathStringForBookBuildingSpan := ""
	xpathStringForCompanyNameSpan := ""
	xpathStringForTargetCd := ""

	bookBuildingString := ""
	companyNameString := ""
	targetCdString := ""

	var companyNameStringList []string
	var targetCdStringList []string
	var bookBuildingStringList []string
	var bookBuildingPossibleBoolListForSbi []string
	var bookBuildingPossibleBoolListForMizuho []string

	for i := 1; i <= 50; i++ {

		bookBuildingPossibleBoolStringForSbi := "false"
		bookBuildingPossibleBoolStringForMizuho := "false"

		xpathStringForCompanyNameSpan = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[2]", i)
		xpathStringForBookBuildingSpan = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[3]", i)
		xpathStringForTargetCd = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[1]/a", i)

		companyNameString, _ = page.FindByXPath(xpathStringForCompanyNameSpan).Text()
		targetCdString, _ = page.FindByXPath(xpathStringForTargetCd).Text()
		bookBuildingString, err = page.FindByXPath(xpathStringForBookBuildingSpan).Text()
		if err != nil {

			break
		} else {

			companyNameStringList = append(companyNameStringList, companyNameString)
			targetCdStringList = append(targetCdStringList, targetCdString)
			bookBuildingStringList = append(bookBuildingStringList, bookBuildingString)

			if strings.EqualFold(sbiBookBuildingMap[targetCdString], "true") && strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "true") {

				bookBuildingPossibleBoolStringForSbi = "true;"
			}

			bookBuildingPossibleBoolListForSbi = append(bookBuildingPossibleBoolListForSbi, bookBuildingPossibleBoolStringForSbi)

			if strings.EqualFold(mizuhoBookBuildingMap[targetCdString], "true") && strings.EqualFold(general.CheckBookoBuildingPossible(bookBuildingString), "true") {

				bookBuildingPossibleBoolStringForMizuho = "true;"
			}

			bookBuildingPossibleBoolListForMizuho = append(bookBuildingPossibleBoolListForMizuho, bookBuildingPossibleBoolStringForMizuho)
		}
	}

	c.JSON(http.StatusOK, strings.Join(companyNameStringList[:], ",")+"&"+strings.Join(bookBuildingStringList[:], ",")+"&"+strings.Join(bookBuildingPossibleBoolListForSbi[:], ",")+"&"+strings.Join(bookBuildingPossibleBoolListForMizuho[:], ",")+"&"+strings.Join(targetCdStringList[:], ","))

	return
}

func mizuhoBookBuilding(c echo.Context) (err error) {

	log.Printf("mizuhoBookBuilding")

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=1980,1200",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	defer driver.Stop()
	driver.Start()

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// 対象サイトに移動
	page.Navigate("https://netclub.mizuho-sc.com/mnc/login?rt_bn=sc_top_hd_login")

	time.Sleep(3 * time.Second)

	page.FindByXPath("/html/body/header[1]/div/div[1]/div/div/div[2]/ul/li[2]").Click()

	page.FindByXPath("//*[@id='form01']/p/span/input").Click()

	page.Navigate("https://netclub.mizuho-sc.com/mnc/tr/ipopo?6")

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		target, err := page.AllByXPath("/html/body/div[2]/div[4]/div[2]/span[2]/span/table/tbody/tr[" + strconv.Itoa(i) + "]/td[2]/span").Text()
		if err != nil {

			// NOOP
		}

		targetXpathString := "/html/body/div[2]/div[4]/div[2]/span[2]/span/table/tbody/tr[" + strconv.Itoa(i) + "]/td[1]/ul/li/a"

		m[target] = targetXpathString
	}

	tickerSymbol := c.Param("tickerSymbol")

	log.Printf(tickerSymbol)

	page.FindByXPath(m[tickerSymbol]).Click()

	time.Sleep(3 * time.Second)

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/div[2]/div[1]/div/p/input").Click()

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/form/div/p/input[1]").Click()

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/span/form/div/p/input").Click()

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/span[1]/span[2]/span[1]/p/span").Click()

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/span[1]/span[2]/span[2]/p/span").Click()

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/span[1]/span[2]/span[3]/p/span").Click()

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/span[1]/form/div[2]/p/input").Click()

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/div[3]/div[1]/form/table/tbody/tr[4]/td/div/table/tbody/tr[1]/td[2]/span/img").Click()

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/div[3]/div[1]/form/table/tbody/tr[5]/td/table/tbody/tr/td[1]/span/input[1]").Click()

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/div[3]/div[1]/form/div[3]/p/input").Click()

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/div[3]/div[1]/form/table/tbody/tr[8]/td/input").Fill("2137")

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/div[3]/div[1]/form/div[3]/ul/li/div[1]/span/input").Click()

	time.Sleep(5 * time.Second)

	page.FindByXPath("/html/body/div[2]/div[4]/div[2]/div[3]/div[1]/form/div[5]/p/input").Click()

	resultString, _ := page.Title()

	return c.JSON(http.StatusOK, resultString)
}

func sbiBookBuilding(c echo.Context) (err error) {

	log.Printf("sbiBookBuilding")

	driver := agouti.ChromeDriver(

		agouti.ChromeOptions("args", []string{

			"--headless",
			"--window-size=1980,1200",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
		}),
	)

	defer driver.Stop()
	driver.Start()

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// 対象サイトに移動
	page.Navigate("https://www.sbisec.co.jp/ETGate")

	// // ユーザーネーム
	// page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))

	// // パスワード
	// page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))

	// 「ログイン」
	page.FindByXPath("/html/body/table/tbody/tr[1]/td[2]/div[2]/form/p[2]/input").Click()

	// 「ブックビルディング情報」
	page.Navigate("https://site2.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETmgR001Control&_DataStoreID=DSWPLETmgR001Control&burl=search_domestic&dir=ipo%2F&file=stock_info_ipo.html&cat1=domestic&cat2=ipo&getFlg=on&int_pr1=150313_cmn_gnavi:6_dmenu_04")

	// 「新規上場株式ブックビルディング／購入意思表示」
	page.FindByXPath("/html/body/div[4]/div/table/tbody/tr/td[1]/div/div[10]/div/div/a/img").Click()

	time.Sleep(3 * time.Second)

	m := make(map[string]string)

	for i := 0; i < 50; i++ {

		targetCd := ""

		targetXpathTableString := "false"

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

			if target == "" {

				targetXpathTableString = strconv.Itoa(i)
			} else {

				targetXpathTableString = "false"
			}
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

	resultString, _ := page.AllByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/table[3]/tbody/tr[1]/td").Text()

	fmt.Println(resultString)

	return c.JSON(http.StatusOK, resultString)
}

func msgHandler(c echo.Context) (err error) {

	fmt.Printf("msgHandler通ってる")

	websocket.Handler(func(ws *websocket.Conn) {

		defer ws.Close()

		// 初回のメッセージを送信
		err := websocket.Message.Send(ws, "Server: Hello, Client!")
		if err != nil {
			c.Logger().Error(err)
		}

		time.Sleep(5 * time.Second)

		// 初回のメッセージを送信
		err2 := websocket.Message.Send(ws, "Server: Hello, Client2!")
		if err2 != nil {
			c.Logger().Error(err2)
		}

		time.Sleep(5 * time.Second)

		// 初回のメッセージを送信
		err3 := websocket.Message.Send(ws, "Server: Hello, Client3!")
		if err3 != nil {
			c.Logger().Error(err3)
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
