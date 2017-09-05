#!/bin/sh

MYDIR=`dirname $0`
if [ -e pastebin.db ]; then
   echo "Not killing an existing database"
   exit 0
fi

sqlite3 pastebin.db < $MYDIR/gendb.sqlite
