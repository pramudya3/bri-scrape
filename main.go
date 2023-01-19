package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/otiai10/gosseract/v2"
)

const (
	loginUrl = "https://ib.bri.co.id/ib-bri"
	username = ""
	password = ""
)

func captcha2Text(captcha []byte) string {
	client := gosseract.NewClient()
	client.SetImageFromBytes(captcha)
	text, err := client.Text()
	if err != nil {
		log.Fatal("error parse image to text", err)
	}
	return text
}

func chromium() *rod.Page {
	u := launcher.New().Bin("/usr/bin/chromium-browser").MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(loginUrl).MustWindowMaximize()
	return page
}

func edge() *rod.Page {
	u := launcher.New().Bin("/usr/bin/microsoft-edge").MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(loginUrl).MustWindowMaximize()
	return page
}

func chrome() *rod.Page {
	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(loginUrl).MustWindowMaximize()
	return page
}

func main() {
	// Create file
	file, _ := os.Create("saldo.csv")
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Select Browser :
	// page := chromium()
	page := edge()
	// page := chrome()

	// Get captcha image
	captcha, _ := page.MustElement("#simple_img > img").MustWaitVisible().Screenshot(proto.PageCaptureScreenshotFormatPng, 1500)

	// Parse image to text
	text := captcha2Text(captcha)

	// Login
	if len(text) > 4 {
		page.MustElement("#loginForm > div.validation > input[type=text]").MustInput(text[1:5])
	}
	page.MustElement("#loginForm > div.validation > input[type=text]").MustInput(text).WaitVisible()
	page.MustElement("#loginForm > input[type=text]:nth-child(5)").MustInput(username).WaitVisible()
	page.MustElement("#loginForm > input[type=password]:nth-child(8)").MustInput(password).WaitVisible()
	page.MustElement("#loginForm > button").MustClick().WaitInvisible()

	// Homepage after login
	page.MustElement("#myaccounts").MustClick().WaitVisible()

	// Get Saldo Tabungan
	fr1 := page.MustElement("#iframemenu").MustFrame()
	fr1.MustElement("body > div.submenu.active > div:nth-child(2) > a").MustClick().MustWaitVisible()
	time.Sleep(3 * time.Second)
	page.MustScreenshot("total-saldo.png")

	header := []string{"Saldo Tabungan :"}
	writer.Write(header)

	// Get tabel saldo tabungan
	fr2 := page.MustElement("#content").MustFrame()
	noRek := fr2.MustElement("#Any_0 > td:nth-child(1)").MustText()
	jenisProduk := fr2.MustElement("#Any_0 > td:nth-child(2)").MustText()
	nama := fr2.MustElement("#Any_0 > td:nth-child(3)").MustText()
	mataUang := fr2.MustElement("#Any_0 > td:nth-child(4)").MustText()
	saldo := fr2.MustElement("#Any_0 > td:nth-child(5)").MustText()

	fmt.Printf("\nNomor Rekening : %s\n\nJenis Produk : %s\n\nNama : %s\n\nMata Uang : %s\n\nSaldo : %s\n\n", noRek, jenisProduk, nama, mataUang, saldo)

	// Export data to file
	data := [][]string{
		{"Nomor Rekening : " + noRek},
		{"Jenis Produk : " + jenisProduk},
		{"Nama : " + nama},
		{"Mata Uang :" + mataUang},
		{"Saldo : " + saldo},
	}
	for _, row := range data {
		_ = writer.Write(row)
	}

	time.Sleep(1 * time.Second)

	// Logout
	page.MustElement("#main-page > div.headerwrap > div > div.uppernav.col-1-2 > span:nth-child(1) > a:nth-child(4)").MustClick()
}
