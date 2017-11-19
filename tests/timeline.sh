#!/bin/sh

. ./testutils.sh

setuptest

SITE=localhost:9090


echo $SITE
# Because the testing server is local it needs to be tested over http; not generating certs
# for test deployments.
authpost http://$SITE/api/v1/post/ notsecret text "An important treatise on the value of friendship"

#OUT=$(authgetvalue http://$SITE/api/v1/timeline/ notsecret)
authgetvalue http://$SITE/api/v1/timeline/ notsecret

teardowntest
