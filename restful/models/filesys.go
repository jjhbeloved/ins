package models

import (
	"errors"
	"strconv"
	"humbird.org/api/utils"
)

const FS = "local"
const SFTP = "sftp"
const FTP = "ftp"
const SSH = "ssh"

var (
	FSList map[string]*FSBasic
	RSList map[string]*Release
)

func init() {
	FSList = make(map[string]*FSBasic)
	RSList = make(map[string]*Release)
}

/**
Path 		作为临时上传的目录
Version		语法版本	v1, v2
 */
type FSBasic struct {
	Id           string        `json:"id"`
	Protocol     string        `json:"protocol"`
	Version      string        `json:"version"`
	IP           string        `json:"ip"`
	Port         int           `json:"port"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	CertPassword string        `json:"certpassword"`
	CertPath     string        `json:"certpath"`
	Path         string        `json:"path"`
}

func (fs *FSBasic) createId() string {
	return fs.Protocol + "-" + fs.IP + "-" + strconv.Itoa(fs.Port) + "-" + fs.Username
}

func GetFS(id string) (*FSBasic, error) {
	if fs, ok := FSList[id]; ok {
		return fs, nil
	}
	return nil, errors.New(id + "not exists")
}

func GetAllFS() map[string]*FSBasic {
	return FSList
}

func AddFS(fs *FSBasic) string {
	fs.Id = fs.createId()
	FSList[fs.Id] = fs
	return fs.Id
}

func DeleteFS(id string) {
	delete(FSList, id)
}

type Release struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	From    *FSBasic        `json:"from"`
	To      *FSBasic        `json:"to"`
	Filters []string        `json:"filters"`
	Pkgs    []*Pkg          `json:"pkgs"`
}

/**
Path 		作为最终解压后文件的目录
IsBak		是否备份
BackPath	备份目录
TmpPath		下载到本地的临时目录, 结合语法版本使用, 尽量填写
CMD		如果填写了CMD shell语法, 内置语法不会生效
Method		这个包携带的动作 deploy, restart, start, stop
 */
type Pkg struct {
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Profile  string
	Path     string `json:"path"`
	IsBak    bool   `json:"isbak"`
	BackPath string `json:"backpath"`
	TmpPath  string `json:"tmppath"`
	CMD      []string `json:"cmd"`
	Method   string `json:"mtd"`
}

func AddRelease(rs *Release) string {
	rs.Id = utils.SessionId()
	RSList[rs.Id] = rs
	return rs.Id
}

func GetRelease(id string) (*Release, error) {
	if fs, ok := RSList[id]; ok {
		return fs, nil
	}
	return nil, errors.New(id + " not exists")
}

func GetAllRelease() map[string]*Release {
	return RSList
}

func UpdateRelease(rs *Release) string {
	RSList[rs.Id] = rs
	return rs.Id
}

func DeleteRelease(id string) {
	delete(RSList, id)
}

type FileOK struct {
	Name string        `json:"name"`
	Size int64         `json:"size"`
	Time string        `json:"time"`
}

type FromTo struct {
	From []*FileOK        `json:"from"`
	To   []*FileOK        `json:"to"`
}
