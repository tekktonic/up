#!/bin/sh

MYDIR=`dirname $0`
if [ -e up.db ]; then
   echo "Not killing an existing database"
   exit 0
fi

sqlite3 up.db < $MYDIR/gendb.sqlite
