package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//get girl
func main() {
	runIt()
}

func runIt() {
	for i := 1; i <= 199; i++ {
		mm_url := "http://www.meizitu.com/a/list_1_" + strconv.Itoa(i) + ".html"
		doc, err := goquery.NewDocument(mm_url)
		if err != nil {
			panic(err)
			return
		}
		doc.Find(".pic").Each(func(i int, s *goquery.Selection) {
			imgParseUrl, ok := s.Find("a").Attr("href")
			if ok {
				fmt.Println(imgParseUrl)
				parentUrl(imgParseUrl)
			}
		})
		//sleep 1 sec
		time.Sleep(5 * time.Second)
	}
}

func parentUrl(url string) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	FolderArr := strings.Split(url, "/")
	lastFolder := strings.Replace(FolderArr[len(FolderArr)-1], ".html", "", -1)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		time.Sleep(300 * time.Second)
		parentUrl(url)
		return
	}
	doc.Find(".postContent").Find("img").Each(func(i int, s *goquery.Selection) {
		imgParseUrl, ok := s.Attr("src")
		if ok {
			var myPath string
			imgFolderArr := strings.Split(imgParseUrl, "/")
			imgFolder := imgFolderArr[len(imgFolderArr)-2]
			img := imgFolderArr[len(imgFolderArr)-1]

			myPath = filepath.Dir(path) + "\\" + lastFolder + "\\" + imgFolder

			_ = os.MkdirAll(myPath, os.ModePerm)
			myPath = myPath + "\\" + img
			pProcessBar(imgParseUrl, myPath)
			time.Sleep(5 * time.Second)
		}
	})
}

func pProcessBar(strSource string, strDest string) {
	sourceName, destName := strSource, strDest

	// check source
	var source io.Reader
	var sourceSize int64
	if strings.HasPrefix(sourceName, "http://") {
		// open as url
		resp, err := http.Get(sourceName)
		if err != nil {
			fmt.Printf("Can't get %s: %v\n", sourceName, err)
			time.Sleep(300 * time.Second)
			pProcessBar(sourceName, destName)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Server return non-200 status: %v\n", resp.Status)
			return
		}
		i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
		sourceSize = int64(i)
		source = resp.Body
	} else {
		// open as file
		s, err := os.Open(sourceName)
		if err != nil {
			fmt.Printf("Can't open %s: %v\n", sourceName, err)
			return
		}
		defer s.Close()
		// get source size
		sourceStat, err := s.Stat()
		if err != nil {
			fmt.Printf("Can't stat %s: %v\n", sourceName, err)
			return
		}
		sourceSize = sourceStat.Size()
		source = s
	}

	// create dest
	dest, err := os.Create(destName)
	if err != nil {
		fmt.Printf("Can't create %s: %v\n", destName, err)
		return
	}
	defer dest.Close()

	// create bar
	bar := pb.New(int(sourceSize)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	// show percents (by default already true)
	bar.ShowPercent = true

	// show bar (by default already true)
	bar.ShowBar = true

	// no need counters
	bar.ShowCounters = false

	// show "time left"
	bar.ShowTimeLeft = true
	// show "time left"
	bar.ShowTimeLeft = true

	// show average speed
	bar.ShowSpeed = true

	// sets the width of the progress bar
	//bar.SetWidth(80)

	// sets the width of the progress bar, but if terminal size smaller will be ignored
	//bar.SetMaxWidth(80)
	//	bar.Format("<.- >")
	bar.Start()

	// create multi writer
	writer := io.MultiWriter(dest, bar)

	// and copy
	io.Copy(writer, source)
	bar.Finish()
}
