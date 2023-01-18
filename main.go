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

type Saldo struct {
	NoRek       string
	JenisProduk string
	Nama        string
	MataUang    string
	Saldo       string
}

func captcha2Text(captcha []byte) string {
	client := gosseract.NewClient()
	client.SetImageFromBytes(captcha)
	text, err := client.Text()
	if err != nil {
		log.Fatal("error parse image to text", err)
	}
	fmt.Println("captcha text: ", text)
	return text
}

func chromium() *rod.Page {
	u := launcher.New().Bin("/usr/bin/chromium-browser").MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(loginUrl).MustWindowNormal()
	return page
}

func edge() *rod.Page {
	u := launcher.New().Bin("/usr/bin/microsoft-edge").MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(loginUrl).MustWindowMaximize()
	return page
}

func chrome() *rod.Page {
	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(loginUrl).MustWindowNormal()
	return page
}

func main() {
	// browser chromium / chrome / edge
	page := chromium()
	// page := edge()
	// page := chrome()

	// write result to export file
	file, err := os.Create("saldo.csv")
	if err != nil {
		log.Fatalln("error create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// get captcha image
	captcha, _ := page.MustElement("#simple_img > img").MustWaitVisible().Screenshot(proto.PageCaptureScreenshotFormatPng, 1500)

	// parse image to text
	text := captcha2Text(captcha)

	// fill form login
	if len(text) > 4 {
		page.MustElement("#loginForm > div.validation > input[type=text]").MustInput(text[1:5])
	}
	page.MustElement("#loginForm > div.validation > input[type=text]").MustInput(text).WaitVisible()
	page.MustElement("#loginForm > input[type=text]:nth-child(5)").MustInput(username).WaitVisible()
	page.MustElement("#loginForm > input[type=password]:nth-child(8)").MustInput(password).WaitVisible()
	page.MustElement("#loginForm > button").MustClick().WaitInvisible()

	// homepage after login
	page.MustElement("#myaccounts").MustClick().WaitVisible()

	// get iframe element
	fr, err := page.MustElement("#iframemenu").Frame()
	if err != nil {
		fmt.Println(err)
	}
	fr.MustElement("body > div.submenu.active > div:nth-child(2) > a").MustClick().MustWaitVisible().MustScreenshot("total-saldo.png")
	fmt.Println("get frame success")

	// get total saldo
	noRek := fr.MustElement("#Any_0 > td:nth-child(1)")
	jenisProduk := page.MustElement("#Any_0 > td:nth-child(2)")
	nama := page.MustElement("#Any_0 > td:nth-child(3)")
	mataUang := page.MustElement("#Any_0 > td:nth-child(4)")
	saldo := page.MustElement("#Any_0 > td:nth-child(5)")

	// create header for export file
	header := []string{"No Rekening", "Jenis Produk", "Nama", "Mata Uang", "Saldo"}
	writer.Write(header)

	// input data to export file
	res := Saldo{}
	res.NoRek = noRek.MustText()
	res.JenisProduk = jenisProduk.MustText()
	res.Nama = nama.MustText()
	res.MataUang = mataUang.MustText()
	res.Saldo = saldo.MustText()
	row := []string{res.NoRek, res.JenisProduk, res.Nama, res.MataUang, res.Saldo}
	writer.Write(row)

	fmt.Printf("Nomor Rekening : %s\n\nJenisProduk : %s\n\nNama : %s\n\nMata Uang : %s\n\nSaldo : %s\n\n", noRek.MustText(), jenisProduk.MustText(), nama.MustText(), mataUang.MustText(), saldo.MustText())

	page.MustWaitIdle()
	time.Sleep(time.Hour)
}
