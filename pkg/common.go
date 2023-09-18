package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strconv"
)

// md5
func PasswordMD5(passwd, salt string) string {
	h := md5.New()
	// 后面增加一个无意义字符串
	h.Write([]byte(passwd + salt + "@.YnO-"))
	cipherStr := h.Sum(nil)
	result := hex.EncodeToString(cipherStr)
	return result
}

// ToInt64 类型转换，获得int64
func ToInt64(v interface{}) (re int64, err error) {
	switch v.(type) {
	case string:
		re, err = strconv.ParseInt(v.(string), 10, 64)
	case float64:
		re = int64(v.(float64))
	case float32:
		re = int64(v.(float32))
	case int64:
		re = v.(int64)
	case int32:
		re = v.(int64)
	default:
		err = errors.New("不能转换")
	}
	return
}
