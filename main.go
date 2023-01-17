package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"github.com/joho/godotenv"
	"github.com/otiai10/gosseract/v2"
)

const (
	LoginURL = "https://ib.bri.co.id/"
)

type Saldo struct {
	NoRek       string
	JenisProduk string
	Nama        string
	MataUang    string
	Saldo       string
}

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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file couldn't be loaded")
	}
	username := os.Getenv("USERNAME_BRI")
	password := os.Getenv("PASSWORD_BRI")

	// create export file (saldo.csv)
	file, err := os.Create("saldo.csv")
	if err != nil {
		log.Fatalln("error create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// launch default browser (chrome)
	browser := rod.New().Timeout(time.Minute).MustConnect()
	defer browser.MustClose()
	page := stealth.MustPage(browser).MustNavigate(LoginURL).MustWindowNormal()

	// launch gosseract (parse image to text)
	client := gosseract.NewClient()
	defer client.Close()

	// get captcha image
	img, err := page.MustElement(".alignimg").MustWaitLoad().Screenshot(proto.PageCaptureScreenshotFormatPng, 1050)
	if err != nil {
		log.Fatal(err)
	}

	// parse image to text
	text := getCaptcha(img)

	// fill login form
	page.MustElement("#loginForm > div.validation > input[type=text]").MustInput(text)
	page.MustElement("#loginForm > input[type=text]:nth-child(5)").MustInput(username)
	page.MustElement("#loginForm > input[type=password]:nth-child(8)").MustInput(password)
	page.MustElement("#loginForm > button").MustClick()

	// create header for export file
	header := []string{"No Rekening", "Jenis Produk", "Nama", "Mata Uang", "Saldo"}
	writer.Write(header)

	// get total saldo
	page.MustElement("#myaccounts").MustClick()
	page.MustElement("body > div.submenu.active > div:nth-child(2) > a").MustClick().MustWaitLoad().MustScreenshot("get_saldo.png")
	noRek := page.MustElement("#Any_0 > td:nth-child(1)")
	jenisProduk := page.MustElement("#Any_0 > td:nth-child(2)")
	nama := page.MustElement("#Any_0 > td:nth-child(3)")
	mataUang := page.MustElement("#Any_0 > td:nth-child(4)")
	saldo := page.MustElement("#Any_0 > td:nth-child(5)")

	// input data to export file
	res := Saldo{}
	res.NoRek = noRek.MustText()
	res.JenisProduk = jenisProduk.MustText()
	res.Nama = nama.MustText()
	res.MataUang = mataUang.MustText()
	res.Saldo = saldo.MustText()
	row := []string{res.NoRek, res.JenisProduk, res.Nama, res.MataUang, res.Saldo}
	writer.Write(row)

	// print total saldo
	fmt.Printf("Nomor Rekening : %s\n\nJenisProduk : %s\n\nNama : %s\n\nMata Uang : %s\n\nSaldo : %s\n\n", noRek.MustText(), jenisProduk.MustText(), nama.MustText(), mataUang.MustText(), saldo.MustText())

	time.Sleep(time.Hour)
}
