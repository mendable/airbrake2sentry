#!/bin/bash -e

clear

dir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
cd $dir

curl -v -XPOST -H"Content-type: text/xml" --data-binary @valid_2_3_notice.xml http://127.0.0.1:6633/notifier_api/v2/notices

echo
