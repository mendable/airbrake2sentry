#!/bin/bash -e

clear

curl -v -XPOST -H"Content-type: text/xml" -d "" http://127.0.0.1:6633/notifier_api/v2/notices

echo

