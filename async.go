package main

import (
	"time"

	nxcli "github.com/nayarsystems/nxgo"
	"github.com/nayarsystems/nxgo/nxcore"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/ei"
)

func ret(r interface{}, e error, cb []*js.Object) {
	switch len(cb) {
	case 1:
		cb[0].Invoke(r, WrapError(e))

	case 2:
		if e == nil {
			cb[0].Invoke(r)
		} else {
			cb[1].Invoke(WrapError(e))
		}
	}
}

func WrapError(e error) *js.Object {
	if e == nil {
		return nil
	}
	jserr := js.Global.Get("Object").New()
	if err, ok := e.(*nxcore.JsonRpcErr); ok {
		jserr.Set("code", err.Cod)
		jserr.Set("message", err.Mess)
		jserr.Set("data", err.Dat)
	} else {
		jserr.Set("code", 0)
		jserr.Set("message", e.Error())
		jserr.Set("data", nil)
	}
	return jserr
}

func WrapTask(task *nxcore.Task) *js.Object {
	if task == nil {
		return nil
	}
	jstask := js.Global.Get("Object").New()
	jstask.Set("sendResult", func(res interface{}, cb ...*js.Object) {
		go func() {
			r, e := task.SendResult(res)
			ret(r, e, cb)
		}()
	})
	jstask.Set("sendError", func(code int, msg string, data interface{}, cb ...*js.Object) {
		go func() {
			r, e := task.SendError(code, msg, data)
			ret(r, e, cb)
		}()
	})
	jstask.Set("path", task.Path)
	jstask.Set("method", task.Method)
	jstask.Set("params", task.Params)
	jstask.Set("tags", task.Tags)
	return jstask
}

func WrapPipe(pipe *nxcore.Pipe) *js.Object {
	if pipe == nil {
		return nil
	}
	jspipe := js.Global.Get("Object").New()
	jspipe.Set("close", func(cb ...*js.Object) {
		go func() {
			r, e := pipe.Close()
			ret(r, e, cb)
		}()
	})
	jspipe.Set("read", func(max int, timeout float64, cb ...*js.Object) {
		go func() {
			r, e := pipe.Read(max, time.Duration(timeout*float64(time.Second)))
			if e != nil {
				ret(nil, e, cb)
				return
			}
			msgs := make([]ei.M, 0)
			for _, msg := range r.Msgs {
				msgs = append(msgs, ei.M{"count": msg.Count, "msg": msg.Msg})
			}
			result := ei.M{"msgs": msgs, "waiting": r.Waiting, "drops": r.Drops}
			ret(result, e, cb)
		}()
	})
	jspipe.Set("write", func(msg interface{}, cb ...*js.Object) {
		go func() {
			r, e := pipe.Write(msg)
			ret(r, e, cb)
		}()
	})
	jspipe.Set("id", pipe.Id())
	return jspipe
}

func WrapNexusConn(nc *nxcore.NexusConn) *js.Object {
	if nc == nil {
		return nil
	}
	jsnc := js.Global.Get("Object").New()
	jsnc.Set("login", func(user string, pass string, cb ...*js.Object) {
		go func() {
			r, e := nc.Login(user, pass)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("setInactivityTimeout", func(timeout float64, cb ...*js.Object) {
		go func() {
			e := nc.SetInactivityTimeout(time.Duration(timeout * float64(time.Second)))
			cb[0].Invoke(e)
			if e == nil {
				<-nc.GetContext().Done()
				cb[1].Invoke(true)
			}
		}()
	})
	jsnc.Set("taskPush", func(method string, params interface{}, timeout float64, cb ...*js.Object) {
		go func() {
			r, e := nc.TaskPush(method, params, time.Duration(timeout*float64(time.Second)))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("taskPull", func(prefix string, timeout float64, cb ...*js.Object) {
		go func() {
			r, e := nc.TaskPull(prefix, time.Duration(timeout*float64(time.Second)))
			ret(WrapTask(r), e, cb)
		}()
	})
	jsnc.Set("taskList", func(prefix string, limit int, skip int, cb ...*js.Object) {
		go func() {
			r, e := nc.TaskList(prefix, limit, skip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("taskListWithOpts", func(prefix string, limit int, skip int, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.TaskList(prefix, limit, skip, getListOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("taskCount", func(prefix string, cb ...*js.Object) {
		go func() {
			r, e := nc.TaskCount(prefix)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("taskCountWithOpts", func(prefix string, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.TaskCount(prefix, getCountOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userCreate", func(user string, pass string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserCreate(user, pass)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userDelete", func(user string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserDelete(user)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userRename", func(user string, new string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserRename(user, new)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userGetTags", func(user string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserGetTags(user)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userGetEffectiveTags", func(user string, prefix string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserGetEffectiveTags(user, prefix)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userSetTags", func(user string, prefix string, tags map[string]interface{}, cb ...*js.Object) {
		go func() {
			r, e := nc.UserSetTags(user, prefix, tags)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userDelTags", func(user string, prefix string, tags []string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserDelTags(user, prefix, tags)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userSetPass", func(user string, pass string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserSetPass(user, pass)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userList", func(prefix string, limit int, skip int, cb ...*js.Object) {
		go func() {
			r, e := nc.UserList(prefix, limit, skip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userListWithOpts", func(prefix string, limit int, skip int, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.UserList(prefix, limit, skip, getListOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userCount", func(prefix string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserCount(prefix)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userCountWithOpts", func(prefix string, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.UserCount(prefix, getCountOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userAddTemplate", func(user string, template string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserAddTemplate(user, template)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userDelTemplate", func(user string, template string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserDelTemplate(user, template)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userAddWhitelist", func(user string, ip string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserAddWhitelist(user, ip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userDelWhitelist", func(user string, ip string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserDelWhitelist(user, ip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userAddBlacklist", func(user string, ip string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserAddBlacklist(user, ip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userDelBlacklist", func(user string, ip string, cb ...*js.Object) {
		go func() {
			r, e := nc.UserDelBlacklist(user, ip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userSetMaxSessions", func(user string, max int, cb ...*js.Object) {
		go func() {
			r, e := nc.UserSetMaxSessions(user, max)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("userSetDisabled", func(user string, disabled bool, cb ...*js.Object) {
		go func() {
			r, e := nc.UserSetDisabled(user, disabled)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("sessionList", func(prefix string, limit int, skip int, cb ...*js.Object) {
		go func() {
			r, e := nc.SessionList(prefix, limit, skip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("sessionListWithOpts", func(prefix string, limit int, skip int, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.SessionList(prefix, limit, skip, getListOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("sessionCount", func(prefix string, cb ...*js.Object) {
		go func() {
			r, e := nc.SessionCount(prefix)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("sessionCountWithOpts", func(prefix string, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.SessionCount(prefix, getCountOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("sessionKick", func(connId string, cb ...*js.Object) {
		go func() {
			r, e := nc.SessionKick(connId)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("sessionReload", func(connId string, cb ...*js.Object) {
		go func() {
			r, e := nc.SessionReload(connId)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("nodeList", func(limit int, skip int, cb ...*js.Object) {
		go func() {
			r, e := nc.NodeList(limit, skip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("node", func(cb ...*js.Object) {
		go func() {
			r, e := nc.Node()
			ret(r, e, cb)
		}()
	})
	jsnc.Set("pipeCreate", func(jopts ei.M, cb ...*js.Object) {
		go func() {
			opts := &nxcore.PipeOpts{
				Length: ei.N(jopts).M("length").IntZ(),
			}
			r, e := nc.PipeCreate(opts)
			ret(WrapPipe(r), e, cb)
		}()
	})
	jsnc.Set("pipeOpen", func(id string, cb ...*js.Object) {
		go func() {
			r, e := nc.PipeOpen(id)
			ret(WrapPipe(r), e, cb)
		}()
	})
	jsnc.Set("topicPublish", func(topic string, msg interface{}, cb ...*js.Object) {
		go func() {
			r, e := nc.TopicPublish(topic, msg)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("topicSubscribe", func(jspipe *js.Object, topic string, cb ...*js.Object) { // !!!Warning, Don't work!!!
		go func() {
			par := ei.M{
				"pipeid": jspipe.Get("id").String(),
				"topic":  topic,
			}
			r, e := nc.Exec("topic.sub", par)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("topicUnsubscribe", func(jspipe *js.Object, topic string, cb ...*js.Object) {
		go func() {
			par := ei.M{
				"pipeid": jspipe.Get("id").String(),
				"topic":  topic,
			}
			r, e := nc.Exec("topic.unsub", par)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("topicList", func(prefix string, limit int, skip int, cb ...*js.Object) {
		go func() {
			r, e := nc.TopicList(prefix, limit, skip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("topicListWithOpts", func(prefix string, limit int, skip int, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.TopicList(prefix, limit, skip, getListOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("topicCount", func(prefix string, cb ...*js.Object) {
		go func() {
			r, e := nc.TopicCount(prefix)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("topicCountWithOpts", func(prefix string, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.TopicCount(prefix, getCountOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("lock", func(lock string, cb ...*js.Object) {
		go func() {
			r, e := nc.Lock(lock)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("unlock", func(lock string, cb ...*js.Object) {
		go func() {
			r, e := nc.Unlock(lock)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("lockList", func(prefix string, limit int, skip int, cb ...*js.Object) {
		go func() {
			r, e := nc.LockList(prefix, limit, skip)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("lockListWithOpts", func(prefix string, limit int, skip int, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.LockList(prefix, limit, skip, getListOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("lockCount", func(prefix string, cb ...*js.Object) {
		go func() {
			r, e := nc.LockCount(prefix)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("lockCountWithOpts", func(prefix string, opts *js.Object, cb ...*js.Object) {
		go func() {
			r, e := nc.LockCount(prefix, getCountOpts(opts))
			ret(r, e, cb)
		}()
	})
	jsnc.Set("exec", func(method string, params interface{}, cb ...*js.Object) {
		go func() {
			r, e := nc.Exec(method, params)
			ret(r, e, cb)
		}()
	})
	jsnc.Set("close", func(cb ...*js.Object) {
		go func() {
			nc.Close()
			ret(nil, nil, cb)
		}()
	})
	jsnc.Set("closed", nc.Closed)
	jsnc.Set("version", func() string { return nxcli.Version.String() })
	jsnc.Set("nexusVersion", func() string { return nc.NexusVersion.String() })
	jsnc.Set("ping", func(timeout float64, cb ...*js.Object) {
		go func() {
			e := nc.Ping(time.Duration(timeout * float64(time.Second)))
			ret(nil, e, cb)
		}()
	})
	return jsnc
}
