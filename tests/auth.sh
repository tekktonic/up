#!/bin/sh

. ./testutils.sh


setuptest
noauthpost http://$SITE/api/v1/post/
check Post accepting unauthenticated calls.

badauthpost http://$SITE/api/v1/post/ neverapassword text foo 
check Post is accepting badly authenticated calls

noauthget http://$SITE/api/v1/timeline
check Timeline is accepting unauthenticated calls

badauthget http://$SITE/api/v1/post?text=foo neverapassword
check Timeline is accepting badly authenticated calls

teardowntest
echo "All auth tests passed."
