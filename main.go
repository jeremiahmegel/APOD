package main

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
)

const urlBase = "https://apod.nasa.gov/apod/"
const htmlPath = "astropix.html"

func downloadImage(imageURL string, fileLocation string) {
	imgResp, _ := http.Get(imageURL)
	defer imgResp.Body.Close()
	img, _ := ioutil.ReadAll(imgResp.Body)
	ioutil.WriteFile(fileLocation, img, 0644)
}

func getImagePath() string {
	resp, _ := http.Get(urlBase + htmlPath)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	imgRegex, _ := regexp.Compile("<a href=\"(image/[^\"]+)\"")
	imgPath := imgRegex.FindSubmatch(body)
	return string(imgPath[1])
}

func main() {
	imgPath := getImagePath()
	downloadImage(urlBase + string(imgPath), "/home/jm/Pictures/Wallpapers/apod.jpg")
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file:///home/jm/Pictures/Wallpapers/apod.jpg")
	cmd.Start()
}
