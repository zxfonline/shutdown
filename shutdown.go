// Copyright 2016 zxfonline@sina.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package shutdown

import (
	"container/list"

	"github.com/zxfonline/golog"
)

var (
	logger *golog.Logger = golog.New("ShutdownHooker")
)

type StopNotifier interface {
	Close()
}
type ShutdownHooker struct {
	hookList *list.List
}

var Hooker = &ShutdownHooker{
	hookList: list.New(),
}

func (this *ShutdownHooker) RegistHook(hook StopNotifier) StopNotifier {
	this.hookList.PushBack(hook)
	return hook
}
func (this *ShutdownHooker) Close() {
	for e := this.hookList.Front(); e != nil; e = e.Next() {
		v := e.Value.(StopNotifier)
		shutdown(v)
	}
	this.hookList.Init()
}
func shutdown(hook StopNotifier) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("recover error:%v, hook=%+v", e, hook)
		} else {
			logger.Debugf("closed hook=%+v", hook)
		}
	}()
	hook.Close()
}

func RegistHook(hook StopNotifier) StopNotifier {
	return Hooker.RegistHook(hook)
}
func Close() {
	Hooker.Close()
}
