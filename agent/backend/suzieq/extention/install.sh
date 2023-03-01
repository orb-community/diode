#!/bin/sh

BASEDIR=$(dirname "$0")
SQPATH=$(pip show suzieq |  sed -n 's/Location: \(.*\)/\1/p')

if [ -z "$SQPATH" ]; then
    exit 1
fi

cp -r "$BASEDIR/poller/" "$SQPATH/suzieq/"
cp -r "$BASEDIR/db/" "$SQPATH/suzieq/"

echo "Installed logging db extention successfully"

exit 0

