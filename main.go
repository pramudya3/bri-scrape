package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"github.com/otiai10/gosseract/v2"
)

const (
	LoginURL = "https://ib.bri.co.id/"
)

func getCaptcha(img []byte) string {
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
	browser := rod.New().Timeout(time.Minute).MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)
	page.MustNavigate(LoginURL).MustWindowNormal()

	// get captcha text
	client := gosseract.NewClient()
	defer client.Close()

	img, err := page.MustElement(".alignimg").MustWaitLoad().Screenshot(proto.PageCaptureScreenshotFormatPng, 1000)
	if err != nil {
		log.Fatal(err)
	}
	text := getCaptcha(img)

	// isi form login
	page.MustElement("#loginForm > div.validation > input[type=text]").MustInput(text)
	page.MustElement("#loginForm > input[type=text]:nth-child(5)").MustInput("")
	page.MustElement("#loginForm > input[type=password]:nth-child(8)").MustInput("")
	page.MustElement("#loginForm > button").MustClick().GetSessionID()

	// get total rekening
	page.MustElement("#myaccounts > table").MustClick()
	page.MustElement("body > div.submenu.active > div:nth-child(2) > a").MustClick()
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
