package main

import (
	"log"
	"os"
	"time"
)

// GetAllFile 递归获取指定目录下的所有文件名
func GetAllFile(pathname string) ([]string, error) {
	var result []string

	fis, err := os.ReadDir(pathname)
	if err != nil {
		log.Fatalf("读取文件目录失败，pathname=%v, err=%v \n", pathname, err)
		return result, err
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp, err := GetAllFile(fullname)
			if err != nil {
				log.Fatalf("读取文件目录失败,fullname=%v, err=%v \n", fullname, err)
				return result, err
			}
			result = append(result, temp...)
		} else {
			result = append(result, fullname)
		}
	}

	return result, nil
}

// SecondToTime 把秒级的时间戳转为time格式
func SecondToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}
