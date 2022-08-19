package main

import (
	"bufio"
	"errors"
	"flag"
	"github.com/go-co-op/gocron"
	"io"
	"log"
	"os"
	path2 "path"
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

func main() {

	path, err2 := os.Open("./path.txt")

	if err2 != nil {
		log.Fatalln("读取path.txt失败")
	}
	path.Close()

	flag.Parse()

	if *rightNow {
		task[runtime.GOOS]()
		return
	}

	scheduler := gocron.NewScheduler(time.Local)

	log.Println("开启自动清理日志, 间隔天数为: ", *days)

	_, err := scheduler.Every(*days).Days().At("0:00").Do(task[runtime.GOOS])
	if err != nil {
		log.Fatalln(err)
		return
	}
	scheduler.StartBlocking()
}

func unixTask() {

	path, err2 := os.Open("./path.txt")

	defer path.Close()
	if err2 != nil {
		log.Fatalln("path.txt文件打开失败")
	}

	reader := bufio.NewReader(path)

	for {
		line, _, err2 := reader.ReadLine()
		if err2 == io.EOF {
			break
		}
		err := unixDeleteFile(string(line))
		if err != nil {
			log.Fatalln(err)
			return
		} else {
			log.Println("文件夹: ", string(line), strconv.Itoa(*days), "天前日志已清除")
		}
	}

}

func unixDeleteFile(path string) error {

	pathFile, err := os.Stat(path)

	if err != nil || !pathFile.IsDir() {
		return errors.New("路径" + path + "无法打开或者非文件夹")
	}

	overDay := time.Now().AddDate(0, 0, -(*days))

	files, _ := GetAllFile(path)

	for i := 0; i < len(files); i++ {

		findFile, err := os.Open(files[i])

		if err != nil {
			findFile.Close()
			log.Fatalln(err)
		}

		stat, err := findFile.Stat()

		if err != nil {
			findFile.Close()
			log.Fatalln(err)
		}

		fileAttr := stat.ModTime()

		if fileAttr.Before(overDay) && path2.Ext(files[i]) == ".log" {
			err := os.Remove(files[i])

			log.Println("删除文件: ", files[i])

			if err != nil {
				findFile.Close()
				return err
			}

			findFile.Close()
		}
	}
	return nil
}

func windowTask() {
	log.Println("暂未支持Windows系统")
	//nothing to do
}
