#!/bin/sh

SCRIPTSDIR=`dirname $0`

ECHO=`which echo`
$ECHO "Time to get things up and running"

$ECHO "What's your name?"
read NAME
$ECHO "What's your domain?"
read DOMAIN
$ECHO "What's your port?"
read PORT

$ECHO "What's your shared secret?"
read SECRET

$ECHO "What should the database file be named?"
read $DBFILE

$ECHO "What should the config file be named?"
read CONFIG


# Generate the config file.
(
cat <<EOF
{
    "owner" : "$NAME",
    "domain" : "$DOMAIN",
    "key" : "$SECRET",
    "max" : 1000,
    "dbfile" : "$DBFILE",
    "port" : "$PORT",
    "timelinesize" : 50,
}
EOF
) > $CONFIG

# Grab dependencies and generate our database.
$SCRIPTSDIR/dependencies.sh && $SCRIPTSDIR/gendb.sh $DBFILE
