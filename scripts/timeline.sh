#!/bin/sh

. `dirname $0`/config.sh

curl -q -H "X-Up-Auth: $SECRET" $SITE/api/v1/timeline/
