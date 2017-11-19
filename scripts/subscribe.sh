#!/bin/sh

. `dirname $0`/config.sh

curl -q -H "X-Up-Auth: $SECRET" -d "$1" $SITE/api/v1/follow/


