package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Download struct {
	URL        string // url for file to download
	TargetPath string // path to download file to
	Sections   int // amount of sections to run concurrently
}

func main() {
	d := Download{
		URL:        "https://download131.uploadhaven.com/1/application/zip/7k8pKJxN5nlMBKHLzielAw4c0CJUdJBrGWy1zMVg.zip?key=Ukd2wo84cUETc08h2uMkqw&expire=1649958272&filename=God.of.War.v1.0.10.zip",
		TargetPath: "./file.zip",
		Sections:   20,
	}

	if err := d.Do(); err != nil {
		fmt.Printf("error occured downloading file: %v\n", err)
		return
	}
	
	fmt.Println("Download successful")
}

func (d *Download) Do() error {
	req, err := d.requestFile("HEAD")
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("response invalid with code %v", resp.StatusCode)
	}

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	fmt.Println("Size of content in bytes:", contentLength)

	// sections is a multi dimensional slice with the outer dimension being the amount of divisions of the file
	// and the inner diment beng an array holding the start and end of a section
	// for example a file of 60 bytes using six sections will have
	// [[0 9][10 19][20 29][30 39][40 49][50 59]]
	sections := make([][2]int, d.Sections)
	eachDivisonSize := contentLength / d.Sections

	// algorithm to split bytes
	for i := range sections {
		// starting bytes
		if i == 0 {
			// first section starting byte will be zero
			sections[i][0] = 0
		} else {
			// other sections starting byte will be the ending byte of previous section + 1
			sections[i][0] = sections[i-1][1] + 1
		}

		// ending bytes
		if i < d.Sections-1 {
			// ending byte of other sections
			sections[i][1] = sections[i][0] + eachDivisonSize
		} else {
			// ending byte of last section
			sections[i][1] = contentLength - 1
		}
	}
	
	var wg sync.WaitGroup
	for idx, section := range sections {
		wg.Add(1)
		go func(i int, s [2]int){
			defer wg.Done()
			if err := d.downloadSection(i, s); err != nil {
				panic(err)
			}
		}(idx, section)
	}
	wg.Wait()

	return nil
}

func (d Download) requestFile(method string) (*http.Request, error) {
	r, err := http.NewRequest(method, d.URL, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", "My concurrent download manager")
	return r, nil
}

func (d Download) downloadSection(idx int, section [2]int) error {
	req, err := d.requestFile("GET")
	if err != nil {
		return err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", section[0], section[1]))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("gotten %v bytes from section %v range %v\n", resp.Header.Get("Content-Length"), idx, section)
	
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(d.TargetPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	written, err := f.Write(bs)
	if err != nil {
		return err
	}

	fmt.Println("Bytes written:", written)

	return nil
}
