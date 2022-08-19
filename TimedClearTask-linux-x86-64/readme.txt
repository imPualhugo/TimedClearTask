输入./clearTask start 14
意为开启自动任务清理,默认为14天

输入./clearTask now
意为立即执行清理程序并结束, 本命令不会开启自动任务

输入./clearTask stop 关闭自动任务
执行清理程序会生成clear.log文件

输入./clearTask find 可以查询清理程序的pid

ps: 不建议使用kill -9 pid杀死进程, 建议使用clearTask stop 或者 pkill -f TimedClearTask关闭进程