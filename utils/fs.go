package utils

import (
	"net/http"
	"os"
	"io"
	"path/filepath"
	"strings"
	"fmt"
)

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		buf = make([]byte, 32 * 1024)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func DownloadToFile(url, filename string) (string, error) {
	dir, _ := filepath.Split(filename)
	os.MkdirAll(dir, os.ModePerm)

	urls := strings.Split(url, ":")
	var err error
	var target string

	switch strings.ToLower(urls[0]) {
	case "http", "https":
		_, err = httpDonwload(url, filename)
		target = filename
		break
	case "file":
		target = urls[1]
		break
	default:
		target = url
		break
	}

	return target, err
}

func DownloadToDir(url, dir string) (string, error) {
	urls := strings.Split(url, ":")
	os.MkdirAll(dir, os.ModePerm)

	var err error
	var target string

	switch strings.ToLower(urls[0]) {
	case "http", "https":
		_, filename := filepath.Split(urls[len(urls) - 1])
		target = filepath.Join(dir, filename)
		_, err = httpDonwload(url, target)
		break
	case "file":
		target = urls[1]
		break
	default:
		target = url
		break
	}
	return target, err
}

func httpDonwload(url, file string) (int64, error) {

	info, err := os.Stat(file)
	if err == nil && info.Size() > 0  {
		return info.Size(), err
	}

	out, err := os.Create(file)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	n, err := copyBuffer(out, resp.Body, nil)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func TempFileserver() {
	http.Handle("/humbird/download/", http.StripPrefix("/humbird/download/", http.FileServer(http.Dir("/veris/odc/install/src"))))
	http.ListenAndServe(":8888", nil)
}

func main() {
	//DownloadToDir("http://10.6.0.13:8888/humbird/download/ruby-2.3.1.tar.gz", "/tmp")
	target, err := DownloadToFile("/tmp/ruby-2.3.1.tar.gz", "/tmp/c.tar.gz")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(target)
}