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
	LoginURL = "https://ib.bri.co.id/login/"
)

type Config struct {
	Username string
	Password string
}

func main() {
	browser := rod.New().Timeout(time.Minute).MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)
	// page := browser.MustPage("https://ib.bri.co.id/ib-bri/login/").MustWindowNormal()
	page.MustNavigate("https://ib.bri.co.id/")

	// get captcha text
	client := gosseract.NewClient()
	defer client.Close()

	page.MustElement("#simple_img > img").MustWaitLoad().MustScreenshot("captcha.png")
	img, _ := page.MustElement("#simple_img > img").MustWaitLoad().Screenshot(proto.PageCaptureScreenshotFormatPng, 400)
	// client.SetImage("captcha.png")
	client.SetImageFromBytes(img)
	text, err := client.Text()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("captcha text: ", text)

	// isi form login
	page.MustElement("#loginForm > div.validation > input[type=text]").MustInput(text)
	page.MustElement("#loginForm > input[type=text]:nth-child(5)").MustInput("")
	page.MustElement("#loginForm > input[type=password]:nth-child(8)").MustInput("")
	page.MustElement("#loginForm > button").MustClick()

	time.Sleep(time.Hour)
}
