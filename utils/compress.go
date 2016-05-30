package utils

import (
	"compress/gzip"
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"
	"path"
	"strings"
	"io"
	"archive/zip"
	"compress/flate"
)

func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func UntarR(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(target, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(target, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
如果使用 archive 就算变更了 .gz的文件名, 也会读取 .gz的时候保存的属性内的文件名
 */
func GunzipL(source string, rm bool) (string, error) {
	reader, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer archive.Close()
	dir, _ := filepath.Abs(filepath.Dir(source))
	fi, _ := reader.Stat()
	target := fi.Name()[0:len(fi.Name()) - 3]
	//dir = filepath.Join(dir, archive.Name)
	dir = filepath.Join(dir, target)

	writer, err := os.Create(dir)
	if err != nil {
		return "", err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	if rm == true && err == nil {
		os.Remove(source)
	}
	return dir, err
}

// Ungzip and untar from source file to destination directory
// you need check file exist before you call this function
func UnTarGz(srcFilePath string, destDirPath string, rm, debug bool) (string, error) {
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)
	fr, err := os.Open(srcFilePath)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	// Gzip reader
	gr, err := gzip.NewReader(fr)
	// Tar reader
	tr := tar.NewReader(gr)
	i := 0
	var name string

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		if debug {
			fmt.Println("UnTarGzing file..." + hdr.Name)
		}
		// Check if it is diretory or file
		if hdr.Typeflag != tar.TypeDir {
			// Get files from archive
			// Create diretory before create file
			os.MkdirAll(destDirPath + "/" + path.Dir(hdr.Name), os.ModePerm)
			// Write data to file
			fw, err := os.Create(destDirPath + "/" + hdr.Name)
			if err != nil {
				return "", err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return "", err
			}
		}
		if i == 0 {
			names := strings.Split(hdr.Name, string(os.PathSeparator))
			if len(names) > 0 {
				name = names[0]
			}
		}
		i++
	}
	if rm {
		os.Remove(srcFilePath)
	}
	return name, nil
}

// 参数frm可以是文件或目录，不会给dst添加.zip扩展名
func Zip(frm, dst string) error {
	zipfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	myzip := zip.NewWriter(zipfile)        // 用压缩器包装该缓冲
	myzip.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	defer myzip.Close()
	err = filepath.Walk(frm,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return filepath.SkipDir
			}
			header, err := zip.FileInfoHeader(info) // 转换为zip格式的文件信息
			if err != nil {
				return filepath.SkipDir
			}
			header.Name, _ = filepath.Rel(filepath.Dir(frm), path)
			w, err := myzip.CreateHeader(header) // 创建一条记录并写入文件信息
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// 确定采用的压缩算法（这个是内建注册的deflate）
				//header.Method = zip.Deflate
			} else {
				header.Name += string(os.PathSeparator)
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(w, file)
			return nil
		})
	return err
}

func Unzip(frm, dst string, rm bool) (string, error) {
	reader, err := zip.OpenReader(frm)
	reader.RegisterDecompressor(zip.Deflate, func(in io.Reader) io.ReadCloser {
		return flate.NewReader(in)
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	if err := os.MkdirAll(dst, 0755); err != nil {
		return "", err
	}
	var name string
	i := 0
	for _, file := range reader.File {
		path := filepath.Join(dst, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			if i == 0 {
				names := strings.Split(file.Name, string(os.PathSeparator))
				if len(names) > 0 {
					name = names[0]
					i++
				}
			}
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			if fileReader != nil {
				fileReader.Close()
			}
			return "", err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			if targetFile != nil {
				targetFile.Close()
			}
			return "", err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			fileReader.Close()
			targetFile.Close()
			return "", err
		}
		fileReader.Close()
		targetFile.Close()
	}

	if rm {
		os.Remove(frm)
	}

	return name, nil
}