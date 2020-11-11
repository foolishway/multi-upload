package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
)

const (
	serverAddr = "10.26.235.42:65345"
)

type uploader struct {
	filePath   string
	size       int64
	wg         *sync.WaitGroup
	serverAddr string
}

func (c *uploader) upload() {
	log.Println(fmt.Sprintf("Uploading %s", c.filePath))
	defer c.wg.Done()

	f, err := os.Open(c.filePath)
	defer f.Close()
	if err != nil {
		panic(fmt.Sprintf("Open file error %v", err))
	}

	conn, err := net.Dial("tcp", serverAddr)
	// defer conn.Close()
	if err != nil {
		panic(fmt.Sprintf("Dial %s error %v", serverAddr, err))
	}

	fLine := fmt.Sprintf("UPLOAD_FILE %s\n", filepath.Base(c.filePath))
	conn.Write([]byte(fLine))
	io.Copy(conn, f)
	log.Println(fmt.Sprintf("Uploaded %s", c.filePath))
}

func main() {
	serverAddr := flag.String("s", "", "Server addr")
	flag.Parse()

	if *serverAddr == "" {
		flag.Usage()
		panic("Server address is required.")
	}

	var files []string
	if flag.NArg() == 0 {
		// panic("Files to upload is requred.")
		finfo, _ := ioutil.ReadDir("./testData")
		for _, fi := range finfo {
			if files == nil {
				files = make([]string, 0)
			}
			files = append(files, "./testData/"+fi.Name())
		}
	} else {
		files = flag.Args()
	}

	var wg *sync.WaitGroup = &sync.WaitGroup{}
	for _, filePath := range files {
		wg.Add(1)
		go func(filePath string) {
			c := uploader{filePath: filePath, wg: wg, serverAddr: *serverAddr}
			c.upload()
		}(filePath)
	}
	wg.Wait()
	log.Println("All files upload completed.")
}