kill $(netstat -ano |grep 9000|grep 0.0.0.0|awk  '{print $5}')
kill $(netstat -ano |grep 5555|grep 0.0.0.0|awk  '{print $5}')
kill $(netstat -ano |grep 5556|grep 0.0.0.0|awk  '{print $5}')
