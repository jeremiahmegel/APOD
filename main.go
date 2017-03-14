package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
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
	forceUpdate := flag.Bool("force", false, "Check for updates even if file has been modified today")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	fileLocation := flag.Arg(0)
	fileLocationAbs, _ := filepath.Abs(fileLocation)

	if !*forceUpdate {
		fileStat, err := os.Stat(fileLocationAbs)
		if err == nil {
			modTime := fileStat.ModTime()
			now := time.Now()
			todayMidnight := now.Truncate(24 * time.Hour)
			if modTime.After(todayMidnight) {
				fmt.Fprintln(os.Stderr, "File has been updated today; not checking for updates")
				os.Exit(0)
			}
		}
	}

	imgPath := getImagePath(urlBase + htmlPath)
	imgFile := downloadImage(urlBase + imgPath)
	os.Rename(imgFile.Name(), fileLocationAbs)
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+fileLocationAbs)
	cmd.Start()
}
