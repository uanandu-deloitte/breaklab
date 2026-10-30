#!/usr/bin/env bash
SRC=/data/to-backup
DEST=/opt/backup-destination
STAMP=$(date +%Y%m%d%H%M%S)
STAGE=/tmp/backup-stage-$STAMP
ARCHIVE=$DEST/backup-$STAMP.tar.gz

mkdir -p $STAGE

for f in $SRC/*; do
  cp $f $STAGE/
done

tar czf $ARCHIVE -C $STAGE .

echo "backup succeeded: $ARCHIVE"
