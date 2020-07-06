#!/usr/bin/env bash
LANG='en_US.UTF-8'
LANGUAGE='en_US.UTF-8'
id -a
for b in `lsblk -l | awk '{print $1}'`; do blkid /dev/$b; done
df