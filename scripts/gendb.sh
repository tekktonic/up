#!/bin/sh

MYDIR=`dirname $0`

if [[ $1 == "" ]]; then
	echo "No db given"
	exit 1
fi

if [ -e $1 ]; then
   echo "Not killing an existing database"
   exit 0
fi

sqlite3 $1 < $MYDIR/gendb.sqlite
