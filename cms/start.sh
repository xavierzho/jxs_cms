#!/bin/sh

/app/data_backend -mode migrate 2>/app/start.log 1>/dev/null
logCnt=$(head -n 1 /app/start.log | wc -l)
if [ $logCnt -gt 0 ]; then
    echo 'migrate fail'
    cat /app/start.log
    exit
fi

/app/data_backend >/dev/null 2>/app/start.log
