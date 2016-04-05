package main

import (
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
)

const (
	url           string = "http://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1"
	tmpImg        string = "temp.jpg"
	backgroundImg string = "background.png"
)

func main() {
	width, height := getScreenResolution()
	res, err := http.Get(url)

	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	imageUrl, err := getInputImage([]byte(body))

	response, err := http.Get(imageUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	usr, err := user.Current()
	err = os.Mkdir(fmt.Sprint(usr.HomeDir, "/.tapet"), 0777)
	if err != nil {
		if !os.IsExist(err) {
			println("dir NOPE")
			log.Fatal(err)
		}
	}

	file, err := os.Create(fmt.Sprint(usr.HomeDir, "/.tapet/", tmpImg))
	if err != nil {
		println("file NOPE")
		log.Fatal(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	colors, err := colorsFromImage(fmt.Sprint(usr.HomeDir, "/.tapet/", tmpImg))

	if len(colors) > 16 {
		colors = colors[:16]
	} else if len(colors) < 16 {
		log.Fatal("Less than 16 colors. Aborting.")
	}

	file, err = os.OpenFile(fmt.Sprint(usr.HomeDir, "/.tapet/", backgroundImg), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	img := randomImage(colors, width, height)
	png.Encode(file, img)
	file.Close()
	changeDesktopBackground(fmt.Sprint(usr.HomeDir, "/.tapet/", backgroundImg))
}
