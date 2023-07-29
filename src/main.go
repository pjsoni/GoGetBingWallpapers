package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Image struct {
	URL  string `xml:"url"`
	Date string `xml:"startdate"`
}

type Images struct {
	XMLName xml.Name `xml:"images"`
	Image   []Image  `xml:"image"`
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "bing.conf", "Path to configuration file")
	flag.Parse()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Configuration file not found. Please specify a valid configuration file using the -config flag.\r\n")
		flag.Usage()
		os.Exit(1)
	}

	config, err := readConfig(configFile)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	fmt.Println("Start time:", start.Format(time.RFC3339))

	//response, e := http.Get(imageURL)
	//fileName := "BingWallpaper-" + image.Date + ".jpg"
	downloadImages(config)

	finish := time.Now()
	fmt.Println("Finish time:", finish.Format(time.RFC3339))
	fmt.Println("Total time:", finish.Sub(start))

}

func populateDownloadData(config Config) []Image {
	imageList := []Image{}
	fmt.Println("Building image list...")
	for idx := config.StartIdx; idx <= config.EndIdx; idx++ {
		url := fmt.Sprintf("http://www.bing.com/HPImageArchive.aspx?format=json&idx=%d&n=%d&mkt=%s", idx, config.NumImages, config.Market)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var images Images
		xml.Unmarshal(body, &images)

		for _, image := range images.Image {
			if !contains(imageList, image) {
				imageList = append(imageList, image)
			}
		}
	}
	return imageList
}

func downloadImages(config Config) {
	imageList := populateDownloadData(config)

	err := os.MkdirAll(config.DestinationDir, 0755)
	if err != nil {
		panic(err)
	}

	for _, image := range imageList {
		fmt.Println("Downloading file: ", image.URL)
		imageURL := strings.Replace(image.URL, config.OldResolution, config.NewResolution, 1)
		date := formatDate(image.Date)
		fileName := config.Prefix + date + ".jpg"
		filePath := filepath.Join(config.DestinationDir, fileName)
		if _, err := os.Stat(filePath); err == nil {
			fmt.Println("Image exists: ", fileName)
			continue
		}

		response, e := http.Get(fmt.Sprintf("https://bing.com%s", imageURL))
		if e != nil {
			panic(e)
		}
		defer response.Body.Close()

		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		io.Copy(file, response.Body)
		fmt.Println("Saved image as: ", fileName)
	}
	fmt.Println("Download complete!")
}

func contains(images []Image, image Image) bool {
	for _, img := range images {
		if img.URL == image.URL && img.Date == image.Date {
			return true
		}
	}
	return false
}
func formatDate(date string) string {
	if len(date) != 8 {
		return date
	}
	return date[:4] + "-" + date[4:6] + "-" + date[6:]
}

type Config struct {
	StartIdx       int
	EndIdx         int
	NumImages      int
	Market         string
	Prefix         string
	DestinationDir string
	OldResolution  string
	NewResolution  string
}

func readConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "StartIdx":
			config.StartIdx, err = strconv.Atoi(value)
			if err != nil {
				return Config{}, err
			}
		case "EndIdx":
			config.EndIdx, err = strconv.Atoi(value)
			if err != nil {
				return Config{}, err
			}
		case "NumImages":
			config.NumImages, err = strconv.Atoi(value)
			if err != nil {
				return Config{}, err
			}
		case "Market":
			config.Market = value
		case "Prefix":
			config.Prefix = value
		case "DestinationDir":
			config.DestinationDir = value
		case "OldResolution":
			config.OldResolution = value
		case "NewResolution":
			config.NewResolution = value
		default:
			continue
		}
	}
	if scanner.Err() != nil {
		return Config{}, scanner.Err()
	}
	return config, nil
}
