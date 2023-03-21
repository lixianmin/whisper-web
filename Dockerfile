FROM lixianmin/golang-mini:latest

WORKDIR /app
# 第一个"."是指docker-compose.yml中指定的build-context中指定的目录；第二个"."指WORKDIR
COPY . .

# RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN go env -w GOPRIVATE=*.gitlab.com,*.gitee.com
# RUN go env -w GOSUMDB=off

#########################分界线：之上是通用配置，之下是需要修改的地方######################
# Expose application ports
EXPOSE 8001 8002
ENV APP_NAME=whisper-web

#########################分界线：之上是需要修改的地方，之下是通用配置######################
ENV IMPORT_PATH=github.com/lixianmin/gonsole
RUN go build -ldflags "-w -s -X $IMPORT_PATH.GitBranchName=`git rev-parse --abbrev-ref HEAD` -X $IMPORT_PATH.GitCommitId=`git log --pretty=format:\"%h\" -1` -X '$IMPORT_PATH.GitCommitMessage=`git show -s --format=%s`' -X $IMPORT_PATH.GitCommitTime=`git log --date=format:'%Y-%m-%dT%H:%M:%S' --pretty=format:%ad -1` -X $IMPORT_PATH.AppBuildTime=`date +%Y-%m-%dT%H:%M:%S`" \
    -mod vendor -gcflags "-N -l" -o $APP_NAME .
CMD ./$APP_NAME