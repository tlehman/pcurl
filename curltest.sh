#!/bin/bash

ifname=http://defini.dev.s3-website-us-west-2.amazonaws.com/sd-v1-4.ckpt 
ofname=/Volumes/RAMDisk/

# download function which takes url and range of bytes as input
download() {
  url=$1
  startindex=$2
  endindex=$3
  curl -X GET -H "range: bytes=$start-$end" -o "$ofname$startindex" $url > /dev/null 2>&1
}

getContentSize() {
  url=$1
  curl -sI $url | grep -i Content-Length | awk '{print $2}' | tr -d '\r'
}

# divide by concurrency
concurrency=$1
contentSize=$(getContentSize $ifname)
rangeSize=$((contentSize / concurrency))

for ((i=0; i<$concurrency; i++)) {
  start=$((i * rangeSize))
  end=$((start + rangeSize - 1))
  if [ $i -eq $((concurrency - 1)) ]; then
    end=$contentSize
  fi
  (download $ifname $start $end) &
}

wait
echo "Download complete"
