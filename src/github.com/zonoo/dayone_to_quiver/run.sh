#!/usr/bin/env bash

make
./dayone-to-quiver -v
rm -fr /tmp/foofoofoo.qvnotebook
./dayone-to-quiver -i "/Users/zono/Desktop/exporttech/tech.json" -o "/tmp/dayone.qvnotebook"
