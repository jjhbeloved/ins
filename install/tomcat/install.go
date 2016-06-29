package tomcat

import (
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	TAR_GZ string = "tar.gz"
	ZIP    string = "zip"
)

type Tomcat struct {
	TOMCAT_HOME        string `json:"tomcat_home"`
	Pkg                string `json:"pkg"`
	IsRemove           bool   `json:"isRemove"`
	TOMCAT_NATIVE_HOME string `json:"tomcat_native_home"`
	JAVA_HOME          string `json:"java_home"`
	APR_HOME           string `json:"apr_home"`
	APR_UTIL_HOME      string `json:"apr_util_home"`
	OPENSSL_HOME       string `json:"openssl_home"` // optional
	TOMCAT_NATIVE_PKG  string `json:"tomcat_native_pkg"`
	APR_PKG            string `json:"apr_pkg"`
	APR_UTIL_PKG       string `json:"apr_util_pkg"`
	OPENSSL_PKG        string `json:"openssl_pkg"` // optional
}

func (w *Tomcat) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

func (tomcat *Tomcat) Install() error {
	fileinfo, err := os.Stat(tomcat.Pkg)
	if err != nil {
		return err
	}

	suffix := filepath.Base(fileinfo.Name())
	var file *os.File
	var ezz error
	var cmd *exec.Cmd
	uncompress := filepath.Join(utils.TMPD, "uncompress.sh")
	defer os.Remove(uncompress)
	if strings.HasSuffix(suffix, TAR_GZ) {
		//name, err = utils.UnTarGz(tomcat.Pkg, dir, tomcat.IsRemove, false)
		// copy
		logs.Print(ioutil.WriteFile(
			uncompress,
			[]byte(fmt.Sprintf(tarTomcatPkg,
				tomcat.TOMCAT_HOME,
				tomcat.Pkg, tomcat.TOMCAT_HOME,
			)),
			0750,
		))
	} else if strings.HasSuffix(suffix, ZIP) {
		//name, err = utils.Unzip(tomcat.Pkg, dir, tomcat.IsRemove)
		// copy
		logs.Print(ioutil.WriteFile(
			uncompress,
			[]byte(fmt.Sprintf(unzipTomcatPkg,
				tomcat.TOMCAT_HOME,
				tomcat.Pkg, tomcat.TOMCAT_HOME,
			)),
			0750,
		))
	} else {
		return fmt.Errorf("%s is not exists.", tomcat.Pkg)
	}
	if err != nil {
		return err
	}
	// 根据模板生成domain
	cmd = exec.Command("sh", uncompress)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	_, ezz = os.Stat(tomcat.APR_PKG)
	if ezz != nil {
		goto over
	}
	_, ezz = os.Stat(tomcat.APR_UTIL_PKG)
	if ezz != nil {
		goto over
	}
	_, ezz = os.Stat(tomcat.TOMCAT_NATIVE_PKG)
	if ezz != nil {
		goto over
	}
	_, ezz = os.Stat(tomcat.JAVA_HOME)
	if ezz != nil {
		goto over
	}
	//apr_home, _ = filepath.Abs(filepath.Dir(tomcat.APR_HOME))
	//apr_util_home, _ = filepath.Abs(filepath.Dir(tomcat.APR_UTIL_HOME))
	os.MkdirAll(tomcat.APR_HOME, 0750)
	os.MkdirAll(tomcat.APR_UTIL_HOME, 0750)
	os.MkdirAll(tomcat.TOMCAT_NATIVE_HOME, 0750)
	os.MkdirAll(utils.TMPD, 0750)
	file, _ = utils.TempFile(utils.TMPD, "install_tomcat_", 0750)
	file.Close()
	if len(tomcat.OPENSSL_PKG) <= 0 || len(tomcat.OPENSSL_HOME) <= 0 {
		// start
		logs.Print(ioutil.WriteFile(
			file.Name(),
			[]byte(fmt.Sprintf(tomcatNativeTemplate,
				tomcat.APR_HOME,
				tomcat.APR_PKG, tomcat.APR_HOME,
				tomcat.APR_HOME,
				tomcat.APR_HOME,
				tomcat.APR_UTIL_HOME,
				tomcat.APR_UTIL_PKG, tomcat.APR_UTIL_HOME,
				tomcat.APR_UTIL_HOME,
				tomcat.APR_UTIL_HOME, tomcat.APR_HOME,
				tomcat.APR_UTIL_HOME,
				tomcat.TOMCAT_NATIVE_HOME,
				tomcat.TOMCAT_NATIVE_PKG,
				tomcat.TOMCAT_NATIVE_HOME, tomcat.APR_HOME, tomcat.JAVA_HOME,
			)),
			0750,
		))
	} else {
		// start
		logs.Print(ioutil.WriteFile(
			file.Name(),
			[]byte(fmt.Sprintf(tomcatNativeTemplateWithSSL,
				tomcat.APR_HOME,
				tomcat.APR_PKG, tomcat.APR_HOME,
				tomcat.APR_HOME,
				tomcat.APR_HOME,
				tomcat.APR_HOME,
				tomcat.APR_UTIL_HOME,
				tomcat.APR_UTIL_PKG, tomcat.APR_UTIL_HOME,
				tomcat.APR_UTIL_HOME,
				tomcat.APR_UTIL_HOME, tomcat.APR_HOME,
				tomcat.APR_UTIL_HOME,
				tomcat.OPENSSL_HOME,
				tomcat.OPENSSL_PKG, tomcat.OPENSSL_HOME,
				tomcat.OPENSSL_HOME,
				tomcat.OPENSSL_HOME,
				tomcat.OPENSSL_HOME,
				tomcat.TOMCAT_NATIVE_HOME,
				tomcat.TOMCAT_NATIVE_PKG,
				tomcat.TOMCAT_NATIVE_HOME, tomcat.APR_HOME, tomcat.OPENSSL_HOME, tomcat.JAVA_HOME,
			)),
			0750,
		))
	}
	// 根据模板生成domain
	cmd = exec.Command("sh", file.Name())
	//defer logs.Print(os.Remove(conf))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Remove(file.Name())
		return err
	}
	os.Remove(file.Name())
over:
	return err
}

const tomcatNativeTemplateWithSSL = `#!/bin/bash
mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh configure --prefix=%s
make -j4 && make install
rm -rf %s-SRC

mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh configure --prefix=%s --with-apr=%s
make -j4 && make install
rm -rf %s-SRC

mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh config --prefix=%s -fPIC no-gost no-shared no-zlib
make depend
make install
rm -rf %s-SRC

mkdir -p %s /tmp/tomcat-native-src
tar xzf %s -C /tmp/tomcat-native-src --strip-components=1
cd /tmp/tomcat-native-src/jni/native
sh configure --prefix=%s --with-apr=%s --with-ssl=%s --with-java-home=%s
make -j4 && make install
rm -rf /tmp/tomcat-native-src/
`

const tomcatNativeTemplate = `#!/bin/bash
mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh configure --prefix=%s
make -j4 && make install
rm -rf %s-SRC

mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh configure --prefix=%s --with-apr=%s
make -j4 && make install
rm -rf %s-SRC

mkdir -p %s /tmp/tomcat-native-src
tar xzf %s -C /tmp/tomcat-native-src --strip-components=1
cd /tmp/tomcat-native-src/jni/native
sh configure --prefix=%s --with-apr=%s --with-ssl=yes --with-java-home=%s
make -j4 && make install
rm -rf /tmp/tomcat-native-src/
`

const tarTomcatPkg = `#!/bin/bash
mkdir -p %s
tar xzf %s -C %s --strip-components=1
`

const unzipTomcatPkg = `#!/bin/bash
mkdir -p %s
unzip %s -d %s
`
