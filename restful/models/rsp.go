package models

import (
	"strconv"
)

type RSP struct {
	Code    string        `json:"code"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
}

/**
找不到资源
 */
func HBD_4001(err error) map[string]*RSP {
	rerr := make(map[string]*RSP)
	rerr["error"] = &RSP{
		Code: "HBD_4001",
		Success: false,
		Message: "Couldn't found resource",
		Data:  err.Error(),
	}
	return rerr
}

/**
可编译错误返回
 */
func HBD_xxx(code int, msg string, err error) map[string]*RSP {
	rerr := make(map[string]*RSP)
	rerr["error"] = &RSP{
		Code: "HBD_" + strconv.Itoa(code),
		Success: false,
		Message: "Couldn't found resource",
		Data:  err.Error(),
	}
	return rerr
}

/**
成功 200
 */
func HBD_200(data interface{}) map[string]*RSP {
	success := make(map[string]*RSP)
	success["success"] = &RSP{
		Code: "HBD_200",
		Success: true,
		Message: "Request Success.",
		Data:  data,
	}
	return success
}