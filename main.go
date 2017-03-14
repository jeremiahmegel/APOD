package main

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
)

func downloadImage(imageURL string, fileLocation string) {
	imgResp, _ := http.Get(imageURL)
	defer imgResp.Body.Close()
	img, _ := ioutil.ReadAll(imgResp.Body)
	ioutil.WriteFile(fileLocation, img, 0644)
}

func main() {
	resp, _ := http.Get("https://apod.nasa.gov/apod/astropix.html")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	imgRegex, _ := regexp.Compile("<a href=\"(image/[^\"]+)\"")
	imgPath := imgRegex.FindSubmatch(body)
	downloadImage("https://apod.nasa.gov/" + string(imgPath[1]), "/home/jm/Pictures/Wallpapers/apod.jpg")
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file:///home/jm/Pictures/Wallpapers/apod.jpg")
	cmd.Start()
}
