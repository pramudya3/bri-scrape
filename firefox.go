package main

import (
	"fmt"
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	path, _ := launcher.LookPath()
	// u := launcher.New().Bin("/snap/bin/firefox").MustLaunch()
	u := launcher.New().Bin(path).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage("https://ib.bri.co.id/").MustWindowNormal()

	client := gosseract.NewClient()
	defer client.Close()

	page.MustElement(".alignimg").MustWaitLoad().MustScreenshot("captcha.png")
	client.SetImage("captcha.png")
	text, err := client.Text()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("captcha text: ", text)

}
