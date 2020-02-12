#!/bin/bash

trap 'echo " - Ctrl-C ignored" ' INT
while true ; do
  sleep 30
done

exit 0
