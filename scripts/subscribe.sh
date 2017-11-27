#!/bin/sh

. `dirname $0`/config.sh

curl -q -H "X-Up-Auth: $SECRET" -d "user=$1" $SITE/api/v1/subscribe/


