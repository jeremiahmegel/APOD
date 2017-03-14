package main

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
)

func main() {
	resp, _ := http.Get("https://apod.nasa.gov/apod/astropix.html")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	imgRegex, _ := regexp.Compile("<a href=\"(image/[^\"]+)\"")
	imgPath := imgRegex.FindSubmatch(body)
	imgResp, _ := http.Get("https://apod.nasa.gov/" + string(imgPath[1]))
	img, _ := ioutil.ReadAll(imgResp.Body)
	ioutil.WriteFile("/home/jm/Pictures/Wallpapers/apod.jpg", img, 0644)
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file:///home/jm/Pictures/Wallpapers/apod.jpg")
	cmd.Start()
}
