package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sclevine/agouti"
)

func main() {

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/", sbiBookBuilding)

	// local サーバー
	e.Logger.Fatal(e.Start(":8000"))
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

	page, err := driver.NewPage()
	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// 対象サイトに移動
	page.Navigate("https://www.sbisec.co.jp/ETGate")

	// ユーザーネーム
	page.FindByXPath("//*[@id='user_input']/input").Fill("XXX")

	// パスワード
	page.FindByXPath("//*[@id='password_input']/input").Fill("XXX")

	// 「ログイン」
	page.FindByXPath("//*[@id='SUBAREA01']/form/div/div/div/p[2]/a/input").Click()

	// 「ブックビルディング情報」
	page.Navigate("https://site2.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETmgR001Control&_DataStoreID=DSWPLETmgR001Control&burl=search_domestic&dir=ipo%2F&file=stock_info_ipo.html&cat1=domestic&cat2=ipo&getFlg=on")

	// 「新規上場株式ブックビルディング／購入意思表示」
	page.FindByXPath("/html/body/div[4]/div/table/tbody/tr/td[1]/div/div[10]/div/div/a").Click()

	targets, err := page.FindByXPath("//img[@src='//sbisec.akamaized.net/v3/images/common/trading/b_ipo_moshikomi.gif']").Count()
	if err != nil {

		targets = 0
	}

	var targetsNameList []string

	if targets != 0 {

		for i := 0; i < targets; i++ {

			page.FindByXPath("//img[@src='//sbisec.akamaized.net/v3/images/common/trading/b_ipo_moshikomi.gif']").Click()

			targetName, err := page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[4]/tbody/tr/td/div/font/b").Text()
			if err != nil {

				// NOOP
			} else {

				targetsNameList = append(targetsNameList, targetName)
			}

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[6]/tbody/tr/td/table/tbody/tr/td[1]/table/tbody/tr[2]/td/input").Fill("100")

			page.FindByXPath("//*[@id='strPriceRadio']").Click()

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[8]/tbody/tr/td[1]/table/tbody/tr/td[2]/input").Fill("XXX")

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr[1]/td/form/table[8]/tbody/tr/td[1]/table/tbody/tr/td[3]/input").Click()

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/form/table[7]/tbody/tr[2]/td/input[1]").Click()

			page.FindByXPath("/html/body/table/tbody/tr/td/table[1]/tbody/tr/td/table[1]/tbody/tr/td/table[5]/tbody/tr/td/a").Click()
		}
	}

	c.JSON(http.StatusOK, "対象数："+strconv.Itoa(targets))
	c.JSON(http.StatusOK, "対象："+strings.Join(targetsNameList[:], ","))

	return
}
