package main

import (
	"bufio"
	"flag"
	"github.com/go-co-op/gocron"
	"io"
	"log"
	"os"
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

	flag.Parse()

	if *rightNow {
		task[runtime.GOOS]()
		return
	}

	scheduler := gocron.NewScheduler(time.Local)

	println("开启自动清理日志, 间隔天数为: ", *days)

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
			println("文件夹: ", string(line), strconv.Itoa(*days), "天前日志已清除")
		}
	}

}

func unixDeleteFile(path string) error {
	overDay := time.Now().AddDate(0, 0, *days*-1)

	files, _ := GetAllFile(path)

	for i := 0; i < len(files); i++ {

		fileInfo, err := os.Open(files[i])

		if err != nil {
			fileInfo.Close()
			log.Fatalln(err)
		}

		stat, err := fileInfo.Stat()

		if err != nil {
			fileInfo.Close()
			log.Fatalln(err)
		}

		fileAttr := stat.ModTime()

		if fileAttr.Before(overDay) {
			err := os.Remove(files[i])

			println(time.Now().Format("2006-01-02 15:04:05.000"), "删除文件: ", files[i])

			if err != nil {
				fileInfo.Close()
				return err
			}
		}
	}
	return nil
}

func windowTask() {
	println("暂未支持Windows系统")
	//nothing to do
}
