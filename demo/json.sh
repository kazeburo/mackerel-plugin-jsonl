#!/bin/sh
set -e
export TMPDIR=tmp
mkdir -p tmp
rm -rf tmp/*

rm -f json.log
perl ./json.pl json.log 10
ls -lh json.log
../mackerel-plugin-jsonl --prefix json --log-file json.log -k total.count -j time -a count -k status -j 'status|replace("^(?:([1235])\d{2}|(4)(?:[0-8]\d|9[0-8]))$","${1}${2}xx")|have("2xx","3xx","4xx","499","5xx")' -a group_by_with_percentage -k latency -j latency -a percentile

sleep 1
echo "--------------------"
perl ./json.pl json.log 1200000
ls -lh json.log
time ../mackerel-plugin-jsonl --prefix json --log-file json.log -k total.count -j time -a count -k status -j 'status|replace("^(?:([1235])\d{2}|(4)(?:[0-8]\d|9[0-8]))$","${1}${2}xx")|have("2xx","3xx","4xx","499","5xx")' -a group_by_with_percentage -k latency -j reqtime -a percentile



