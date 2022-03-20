package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sclevine/agouti"
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
	e.GET("/sbiBookBuilding", sbiBookBuilding)
	e.GET("/sbiBalance", getSbiBalance)
	e.GET("/mizuhoBalance", getMizuhoBalance)
	e.GET("/schedule", getSchedule)

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

func gormConnect() *gorm.DB {

	HOST := "db_container"
	PORT := "5432"
	USER := os.Getenv("POSTGRES_USER")
	PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	DBNAME := os.Getenv("POSTGRES_DB")

	CONNECT := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", HOST, PORT, USER, DBNAME, PASSWORD)

	db, err := gorm.Open("postgres", CONNECT)
	if err != nil {

		panic(err.Error())
	}

	return db
}

func sbiBookBuilding(c echo.Context) (err error) {

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

	// ユーザーネーム
	page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))

	// パスワード
	page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))

	// 「ログイン」
	page.FindByXPath("//*[@id='SUBAREA01']/form/div/div/div/p[2]/a/input").Click()

	// 「ブックビルディング情報」
	page.Navigate("https://site2.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETmgR001Control&_DataStoreID=DSWPLETmgR001Control&burl=search_domestic&dir=ipo%2F&file=stock_info_ipo.html&cat1=domestic&cat2=ipo&getFlg=on")

	// 「新規上場株式ブックビルディング／購入意思表示」
	page.FindByXPath("/html/body/div[4]/div/table/tbody/tr/td[1]/div/div[10]/div/div/a").Click()

	targets, err := page.Find("//img[@src='//sbisec.akamaized.net/v3/images/common/trading/b_ipo_moshikomi.gif']").Count()
	if err != nil {

		targets = 0
	}
	fmt.Println("対象取得")

	fmt.Println(targets)

	var targetsNameList []string

	if targets != 0 {

		for i := 0; i < targets; i++ {

			page.FindByXPath("//img[@src='//sbisec.akamaized.net/v3/images/common/trading/b_ipo_moshikomi.gif']").Click()

			targetName, err := page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[4]/tbody/tr/td/div/font/b").Text()
			if err != nil {

				// NOOP
				println(err)
			} else {

				targetsNameList = append(targetsNameList, targetName)
			}

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[6]/tbody/tr/td/table/tbody/tr/td[1]/table/tbody/tr[2]/td/input").Fill("100")

			page.FindByXPath("//*[@id='strPriceRadio']").Click()

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[8]/tbody/tr/td[1]/table/tbody/tr/td[2]/input").Fill(os.Getenv("SBI_TORIHIKI_PASSWORD"))

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[8]/tbody/tr/td[1]/table/tbody/tr/td[3]/input").Click()

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/form/table[7]/tbody/tr[2]/td/input[1]").Click()

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/table[5]/tbody/tr/td/a").Click()
		}
	}

	c.JSON(http.StatusOK, "対象数："+strconv.Itoa(targets))
	c.JSON(http.StatusOK, "対象："+strings.Join(targetsNameList[:], ","))

	return
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

	// ユーザーネーム
	page.FindByXPath("//*[@id='user_input']/input").Fill(os.Getenv("SBI_USERNAME"))

	// パスワード
	page.FindByXPath("//*[@id='password_input']/input").Fill(os.Getenv("SBI_LOGIN_PASSWORD"))

	// 「ログイン」
	page.FindByXPath("//*[@id='SUBAREA01']/form/div/div/div/p[2]/a/input").Click()

	time.Sleep(3 * time.Second)

	page.FindByXPath("/html/body/div[1]/div[1]/div[2]/div/ul/li[3]/a/img").Click()

	time.Sleep(3 * time.Second)

	kaitsukeKano2daysAfter, err := page.FindByXPath("/html/body/div[1]/table/tbody/tr/td[1]/table/tbody/tr[2]/td/table[1]/tbody/tr/td/form/table[2]/tbody/tr[1]/td[2]/table[4]/tbody/tr/td[1]/table[2]/tbody/tr[3]/td[2]/div").Text()
	if err != nil {

		fmt.Printf("err: %v\n", err)
	}
	time.Sleep(3 * time.Second)

	kaitsukeKano3daysAfter, err := page.FindByXPath("/html/body/div[1]/table/tbody/tr/td[1]/table/tbody/tr[2]/td/table[1]/tbody/tr/td/form/table[2]/tbody/tr[1]/td[2]/table[4]/tbody/tr/td[1]/table[2]/tbody/tr[4]/td[2]/div").Text()
	if err != nil {

		fmt.Printf("err: %v\n", err)
	}
	c.JSON(http.StatusOK, "買付余力(2営業日後)"+kaitsukeKano2daysAfter)
	c.JSON(http.StatusOK, "買付余力(3営業日後)"+kaitsukeKano3daysAfter)

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

	page.FindByXPath("//*[@id='IDInputKB']").Fill(os.Getenv("MIZUHO_ID"))

	page.FindByXPath("//*[@id='PWInputKB']").Fill(os.Getenv("MIZUHO_PASSWORD"))

	page.FindByXPath("//*[@id='form01']/p/span/input").Click()

	time.Sleep(5 * time.Second)

	zandaka, err := page.FindByXPath("/html/body/div[2]/div[3]/div/div/div/div[1]/span[2]/div/table/tbody/tr/td[2]/table/tbody/tr/td/span/strong/span").Text()
	if err != nil {

		fmt.Printf("err: %v\n", err)
	}

	c.JSON(http.StatusOK, "買付余力"+zandaka)

	return
}

func getSchedule(c echo.Context) (err error) {

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

	page.Navigate("https://www.nikkei.com/markets/kigyo/ipo/money-schedule/")

	xpathStringForBookBuildingSpan := ""
	xpathStringForCompanyNameSpan := ""

	bookBuildingString := ""
	companyNameString := ""

	// m := make(map[string]string)

	var companyNameStringList []string
	var bookBuildingStringList []string

	for i := 1; i <= 50; i++ {

		xpathStringForCompanyNameSpan = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[2]", i)
		xpathStringForBookBuildingSpan = fmt.Sprintf("/html/body/div[8]/div/div/div/div[3]/div[2]/div[2]/div/div/div[2]/div/table/tbody[1]/tr[%d]/td[3]", i)

		companyNameString, _ = page.FindByXPath(xpathStringForCompanyNameSpan).Text()
		bookBuildingString, err = page.FindByXPath(xpathStringForBookBuildingSpan).Text()
		if err != nil {

			break
		} else {

			companyNameStringList = append(companyNameStringList, companyNameString)
			bookBuildingStringList = append(bookBuildingStringList, bookBuildingString)

			// m[companyNameString] = bookBuildingString
		}
	}

	// data, err := json.Marshal(m)
	// if err != nil {
	// 	fmt.Println("JSON marshal error: ", err)
	// 	return
	// }

	fmt.Println(strings.Join(companyNameStringList[:], ","))
	fmt.Println(strings.Join(bookBuildingStringList[:], ","))

	fmt.Println("test handler update test")

	c.JSON(http.StatusOK, strings.Join(companyNameStringList[:], ",")+"&"+strings.Join(bookBuildingStringList[:], ","))

	return
}
