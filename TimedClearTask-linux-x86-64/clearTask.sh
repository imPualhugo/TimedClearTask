#find /kfts/V8/apache-tomcat-9.0.46-linux/logs/ -mtime +14 -name "*.log" -type f -print -exec rm -rf {} \
#
#echo "V8 14天前日志清理成功"
#
#find /kfts/V8/apache-tomcat-9.0.46-linux/logs/ -mtime +14 -name "*.log" -type f -print -exec rm -rf {} \
#
#echo "trader 14天前日志清理成功"
#
#crontab -e 0 0 0 1/14 *

program=TimedClearTask


if [[ -z "$1" ]]|| [[ "$1" != "now" && "$1" != "start" && "$1" != "stop" && "$1" != "find" ]] ;then
    echo '参数丢失或无效, 具体使用方法参看readme.txt'
    exit 1
fi


if [ "$1" == "now" ]; then
    echo '立即执行日志清理程序'
    if [ -z "$2" ]; then
        ./$program -now
      else
        ./$program -now -day $2
    fi
fi


if [ "$1" == "start" ]; then
    pkill -f $program

    if [ -z "$2" ]; then
         nohup ./$program > clear.log &
         else
           nohup ./$program -day "$2" > clear.log &
    fi
    echo "日志清理程序已开启"
fi

if [ "$1" == "stop" ]; then
    pkill -f $program

    echo "日志清理程序已关闭"
fi

if [ "$1" == "find" ]; then
    pid1=$(pgrep $program)
    if [ $pid1 != '' ]; then
        echo '日志清理程序已启动, pid为'$pid1
        else
          echo '日志清理程序未启动'
    fi
fi