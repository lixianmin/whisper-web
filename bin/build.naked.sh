#!/bin/bash

IMPORT_PATH=github.com/lixianmin/gonsole
FLAGS="-w -s -X $IMPORT_PATH.GitBranchName=`git rev-parse --abbrev-ref HEAD` -X $IMPORT_PATH.GitCommitId=`git log --pretty=format:\"%h\" -1` -X '$IMPORT_PATH.GitCommitMessage=`git show -s --format=%s`' -X $IMPORT_PATH.GitCommitTime=`git log --date=format:'%Y-%m-%dT%H:%M:%S' --pretty=format:%ad -1` -X $IMPORT_PATH.AppBuildTime=`date +%Y-%m-%dT%H:%M:%S`"
go build -ldflags "$FLAGS" -mod vendor -gcflags "-N -l"


# first build, then restart, to avoid the restarted process running the old executable
#APP_NAME=${PWD##*/}
APP_NAME=$(basename `git config --get remote.origin.url`)
APP_NAME=${APP_NAME%.*}  # remove extension if exists

PID=`ps aux | grep ${APP_NAME} | grep -v grep | awk '{ print $2}'`
if [ "$PID" != "" ] ;then
  kill -9 $PID
  echo "kill old process with name="$APP_NAME", pid="$PID
else
  echo "can not find old process with name="$APP_NAME
fi