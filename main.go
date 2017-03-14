package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

const urlBase = "https://apod.nasa.gov/apod/"
const htmlPath = "astropix.html"

func downloadImage(imageURL string) *os.File {
	imgResp, _ := http.Get(imageURL)
	defer imgResp.Body.Close()
	img, _ := ioutil.ReadAll(imgResp.Body)
	tmpFile, _ := ioutil.TempFile("", "apod")
	tmpFile.Write(img)
	return tmpFile
}

func getImagePath(htmlURL string) string {
	resp, _ := http.Get(htmlURL)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	imgRegex, _ := regexp.Compile("<a href=\"(image/[^\"]+)\"")
	imgPath := imgRegex.FindSubmatch(body)
	return string(imgPath[1])
}

func main() {
	args := os.Args
	fileLocation := args[1]
	// TODO Test for correct argument usage

	imgPath := getImagePath(urlBase + htmlPath)
	imgFile := downloadImage(urlBase + imgPath)
	os.Rename(imgFile.Name(), fileLocation)
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+fileLocation)
	cmd.Start()
}
