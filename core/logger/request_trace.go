package logger

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/lixianmin/gonsole/road"
)

/********************************************************************
created:    2020-08-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type RequestTrace struct {
	ctx context.Context
	//requestid 调试用的，请求
	rid       int32
	beginTime time.Time
	session   road.Session
	uid       int64
	marks     []string
}

func NewRequestTrace(ctx context.Context, request interface{}) *RequestTrace {
	var item = &RequestTrace{
		rid:       fetchRequestId(request),
		beginTime: time.Now(),
	}

	if ctx != nil {
		var se = road.GetSessionFromCtx(ctx)
		item.session = se

		if se != nil {
			item.uid = se.Attachment().Int64("uid")
		}
	} else {
		ctx = context.Background()
	}

	var sid = item.session.Id()
	theLogger.Info("sid=%d uid=%d request=%+v", sid, item.uid, request)
	return item
}

func (trace *RequestTrace) Reset() *RequestTrace {
	trace.beginTime = time.Now()
	return trace
}

func (trace *RequestTrace) Info(format interface{}, args ...interface{}) {
	var cost = time.Now().Sub(trace.beginTime)
	var message = formatLog(format, args...)
	// 200ms是个槛：低于200ms人类是无感知的
	// 1s是个槛：超过1s人类会有明显的等待感
	var sid = trace.session.Id()
	if cost < time.Second {
		theLogger.Info("sid=%d uid=%d cost=%dms message=%q", sid, trace.uid, cost.Milliseconds(), message)
	} else {
		var markMessage = strings.Join(trace.marks, ", ")
		message += ", marks=[" + markMessage + "]"
		theLogger.Warn("sid=%d uid=%d cost=%dms message=%q", sid, trace.uid, cost.Milliseconds(), message)
	}
}

func (trace *RequestTrace) NewError(code string, format interface{}, args ...interface{}) *road.Error {
	var cost = time.Now().Sub(trace.beginTime)
	var message = formatLog(format, args...)
	var sid = trace.session.Id()
	theLogger.Info("sid=%d uid=%d cost=%dms code=%q message=%q", sid, trace.uid, cost.Milliseconds(), code, message)

	// 这里必须返回pitaya的Error，否则还会被pitaya重新封装一遍
	var err = road.NewError(code, "sid=%d uid=%d caller=%s message=%q", sid, trace.uid, fetchCallerName(), message)
	return err
}

func (trace *RequestTrace) Mark(name string) {
	var cost = time.Now().Sub(trace.beginTime)
	var mark = fmt.Sprintf("%s:%dms", name, cost.Milliseconds())
	trace.marks = append(trace.marks, mark)
}

func (trace *RequestTrace) GetSession() road.Session {
	return trace.session
}

func (trace *RequestTrace) GetUserId() int64 {
	return trace.uid
}

func (trace *RequestTrace) GetBeginTime() time.Time {
	return trace.beginTime
}

func fetchCallerName() string {
	const skip = 3
	const depth = 1
	var pcs [depth]uintptr // 程序计算器
	var total = runtime.Callers(skip, pcs[:])

	var frames = runtime.CallersFrames(pcs[:total])
	var frame, _ = frames.Next()

	var function = frame.Function
	if function != "" {
		var lastIndex = strings.LastIndexByte(function, '.')
		if lastIndex > 0 {
			var s = function[lastIndex+1:]
			return s
		}
	}

	return ""
}

func formatLog(first interface{}, args ...interface{}) string {
	var msg string
	switch first.(type) {
	case string:
		msg = first.(string)
		if len(args) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(args))
		}
	default:
		msg = fmt.Sprint(first)
		if len(args) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(args))
	}
	return fmt.Sprintf(msg, args...)
}

func fetchRequestId(request interface{}) int32 {
	//获取结构体request的值，原因不知道request是什么类型的，只知道是结构体
	var requestValue = reflect.Indirect(reflect.ValueOf(request))
	var field = requestValue.FieldByName("Id")

	var rid int32 = 0
	if field.Kind() == reflect.Int32 {
		rid = int32(field.Int())
	}

	return rid
}
