package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/otiai10/gosseract/v2"
)

func captcha(img []byte) string {
	client := gosseract.NewClient()
	err := client.SetImageFromBytes(img)
	if err != nil {
		log.Fatalln("error getting captcha.png ", err)
	}
	text, err := client.Text()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("captcha text: ", text)
	return text
}

func main() {
	u := launcher.New().Bin("/snap/bin/firefox")
	page := rod.New().ControlURL(u.MustLaunch()).MustConnect().NoDefaultDevice().MustPage("https://ib.bri.co.id/ib-bri/").MustWindowNormal()

	client := gosseract.NewClient()
	defer client.Close()

	img, err := page.MustElement(".alignimg").MustWaitLoad().Screenshot(proto.PageCaptureScreenshotFormatPng, 1000)
	if err != nil {
		log.Fatal(err)
	}
	text := captcha(img)

	// isi form login
	page.MustElement(".validation > input:nth-child(1)").MustInput(text)
	page.MustElement("#loginForm > input:nth-child(5)").MustInput("rizkypramudya3")
	page.MustElement("#loginForm > input:nth-child(7)").MustInput("Pramudya3")
	page.MustElement("#loginForm > button:nth-child(10)").MustClick().GetSessionID()

	// get total rekening
	page.MustElement("#myaccounts > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1)").MustClick()
	page.MustElement("div.submenu:nth-child(2) > div:nth-child(2) > a:nth-child(1)").MustClick().MustWaitLoad().MustScreenshot("total_saldo.png")
	noRek := page.MustElement("#Any_0 > td:nth-child(1)")
	jenisProduk := page.MustElement("#Any_0 > td:nth-child(2)")
	nama := page.MustElement("#Any_0 > td:nth-child(3)")
	mataUang := page.MustElement("#Any_0 > td:nth-child(4)")
	saldo := page.MustElement("#Any_0 > td:nth-child(5)")

	fmt.Println("Nomor Rekening : ", noRek.MustText())
	fmt.Println("Jeni Produk : ", jenisProduk.MustText())
	fmt.Println("Nama : ", nama.MustText())
	fmt.Println("Mata Uang : ", mataUang.MustText())
	fmt.Println("Saldo : ", saldo.MustText())

	time.Sleep(time.Hour)
}
