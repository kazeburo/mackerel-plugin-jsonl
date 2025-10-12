#!/bin/sh
set -e
export TMPDIR=tmp
mkdir -p tmp
rm -rf tmp/*

rm -f json.log
perl ./json.pl json.log 10
ls -lh json.log
../mackerel-plugin-jsonl --prefix json --log-file json.log -k total.count -j time -a count -k status -j 'status|replace("^(.).+$","${1}xx")' -a group_by_with_percentage -k latency -j latency -a percentile

sleep 1
echo "--------------------"
perl ./json.pl json.log 1200000
ls -lh json.log
time ../mackerel-plugin-jsonl --prefix json --log-file json.log -k total.count -j time -a count -k status -j 'status|replace("^(.).+$","${1}xx")' -a group_by_with_percentage -k latency -j reqtime -a percentile


#ls -lh json.log
#time ../mackerel-plugin-axslog --key-prefix json --logfile json.log --ptime-key=reqtime --format json ---filter example

