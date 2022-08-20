package main

import (
	"bufio"
	"errors"
	"flag"
	"github.com/go-co-op/gocron"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

var task map[string]func()

func init() {

	task = make(map[string]func())

	task["linux"] = unixTask
	task["unix"] = unixTask
	task["darwin"] = unixTask

	task["windows"] = windowTask
}

var days = flag.Int("day", 14, "清理日志天数")

var rightNow = flag.Bool("now", false, "立即执行")

var inputPathFile = flag.String("path", "./path.txt", "读取的路径文件")

func main() {

	flag.Parse()

	pathFile, err2 := os.Open(*inputPathFile)

	if err2 != nil {
		log.Fatalln("读取路径文件: " + (*inputPathFile) + "失败")
		return
	}
	err := pathFile.Close()
	if err != nil {
		pathError(pathFile, err)
		return
	}

	if *rightNow {
		task[runtime.GOOS]()
		return
	}

	scheduler := gocron.NewScheduler(time.Local)

	log.Println("开启自动清理日志, 间隔天数为: ", *days)

	_, err = scheduler.Every(*days).Days().At("0:00").Do(task[runtime.GOOS])
	if err != nil {
		log.Fatalln(err)
		return
	}
	scheduler.StartBlocking()
}

func unixTask() {

	pathFile, err2 := os.Open(*inputPathFile)

	defer pathFile.Close()
	if err2 != nil {
		log.Fatalln("pathFile.txt文件打开失败")
		return
	}

	reader := bufio.NewReader(pathFile)

	for {
		line, _, err2 := reader.ReadLine()
		if err2 == io.EOF {
			break
		}
		err := unixDeleteFile(string(line))
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Println("文件夹: ", string(line), strconv.Itoa(*days), "天前日志已清除")
	}

}

func unixDeleteFile(pathName string) error {

	pathFile, err := os.Stat(pathName)

	overDay := time.Now().AddDate(0, 0, -(*days))

	if err != nil || !pathFile.IsDir() {
		return errors.New("路径" + pathName + "无效")
	}

	//读取到的是文件
	if !pathFile.IsDir() {
		return deleteLogFile(&overDay, pathName)
	}

	files, _ := GetAllFile(pathName)

	for i := 0; i < len(files); i++ {

		err := deleteLogFile(&overDay, files[i])
		if err != nil {
			return err
		}

	}
	return nil
}

// 根据传入的时间和文件路径删除文件
func deleteLogFile(overDay *time.Time, filePath string) error {
	if path.Ext(filePath) != ".log" {
		return nil
	}

	findFile, err := os.Open(filePath)

	if err != nil {
		findFile.Close()
		log.Fatalln(err)
		return err
	}

	stat, err := findFile.Stat()

	if err != nil {
		findFile.Close()
		log.Fatalln(err)
		return err
	}

	fileAttr := stat.ModTime()

	if fileAttr.Before(*overDay) {
		err := os.Remove(filePath)

		log.Println("删除文件: ", filePath)

		if err != nil {
			findFile.Close()
			return err
		}

		findFile.Close()
	}
	return nil
}

func pathError(path *os.File, err error) {
	log.Fatalln("Path: " + path.Name() + "crushed an Error , Error: " + err.Error())
}

func windowTask() {
	log.Println("暂未支持Windows系统")
	//nothing to do
}
