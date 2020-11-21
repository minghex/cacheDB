curl -v 127.0.0.1:12345/cache/testkey -XPUT -dtestvalue

curl -v 127.0.0.1:12345/cache/testkey

curl 127.0.0.1:12345/status

./benchmark -type tcp -n 100000 -r 100000 -t set

./benchmark -type tcp -n 100000 -r 100000 -t get

./benchmark -type tcp -n 100000 -r 100000 -t del