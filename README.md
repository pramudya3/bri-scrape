#bri-scrape


adalah project untuk mendapatkan total saldo dari Bank BRI.
Dengan menggunakan bahasa Golang dibantu dengan library go-rod dan gosseract.
Untuk rate percobaan solver captcha kemungkinan berhasil: 20/2 dan mungkin bisa lebih atau bahkan bisa kurang, tergantung rejeki.


    git clone github.com/pramudya3/bri-scrape
 
 Untuk menjalankan program :
     
    go run main.go -rod=show,slow=1s
    
  Referensi :
  
    https://go-rod.github.io/
    
    https://github.com/otiai10/gosseract
