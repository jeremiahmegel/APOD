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

func mostRecentMidnight(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] filename\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage

	forceUpdate := flag.Bool("force", false, "Check for updates even if file has been modified today")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	fileLocation := flag.Arg(0)
	fileLocationAbs, _ := filepath.Abs(fileLocation)

	if !*forceUpdate {
		fileStat, err := os.Stat(fileLocationAbs)
		if err == nil {
			modTime := fileStat.ModTime()
			now := time.Now()
			midnight := mostRecentMidnight(now)
			if modTime.After(midnight) {
				fmt.Fprintln(os.Stderr, "File has been updated today; not checking for updates")
				os.Exit(1)
			}
		}
	}

	imgPath := getImagePath(urlBase + htmlPath)
	imgFile := downloadImage(urlBase + imgPath)
	os.Rename(imgFile.Name(), fileLocationAbs)
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+fileLocationAbs)
	cmd.Start()
}
