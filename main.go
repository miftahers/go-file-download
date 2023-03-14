package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	url := "https://download.quranicaudio.com/quran/mishaari_raashid_al_3afaasee/"
	startIndex := 1
	endIndex := 2

	err := downloadFiles(url, startIndex, endIndex)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All files downloaded successfully")
}

func downloadFiles(url string, startIndex, endIndex int) error {
	for i := startIndex; i <= endIndex; i++ {
		filename := fmt.Sprintf("%03d.mp3", i)
		fullUrl := url + filename

		fileSize, err := getFileSize(fullUrl)
		if err != nil {
			return fmt.Errorf("error getting file size for %v: %v", fullUrl, err)
		}

		err = downloadFile(fullUrl, filename, fileSize)
		if err != nil {
			return fmt.Errorf("error downloading %v: %v", fullUrl, err)
		}

		fmt.Printf("Downloaded %v\n", filename)
	}

	return nil
}

func getFileSize(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	fileSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return 0, err
	}

	return fileSize, nil
}

func downloadFile(url, filename string, fileSize int) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	var totalBytes int

	for {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			_, err = file.Write(buffer[:n])
			if err != nil {
				return err
			}
			totalBytes += n

			percentage := float64(totalBytes) / float64(fileSize) * 100
			fmt.Printf("\rDownloading %v: %.2f%%", filename, percentage)
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
