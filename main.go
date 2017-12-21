package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/nayarsystems/nxgo"
	"github.com/nayarsystems/nxgo/nxcore"
)

func main() {
	// Node.js
	if require := js.Global.Get("require"); require != js.Undefined {
		if ws := require.Invoke("ws"); ws != js.Undefined {
			println("Running on Node.js with websocket module: 'ws'")

			ws.Get("prototype").Set("removeEventListener", js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
				// Do nothing in nodejs, we can't remove EventListener
				return nil
			}))
			js.Global.Set("WebSocket", ws)
		} else {
			println("Running on Node.js but websocket module missing (try: npm install ws)")
		}
	}

	dial := func(a string, cb func(interface{}, error)) {
		go func() {
			nc, e := nxgo.Dial(a, nil)
			cb(WrapNexusConn(nc), e)
		}()
	}

	var nexus *js.Object

	if js.Module == js.Undefined {
		// Browser
		js.Global.Set("nexus", make(map[interface{}]interface{}))
		nexus = js.Global.Get("nexus")
	} else {
		// Node.js
		nexus = js.Module.Get("exports")
	}

	nexus.Set("dial", dial)
	nexus.Set("errors", nxcore.ErrStr)
	nexus.Set("ErrParse", nxcore.ErrParse)
	nexus.Set("ErrInvalidRequest", nxcore.ErrInvalidRequest)
	nexus.Set("ErrInternal", nxcore.ErrInternal)
	nexus.Set("ErrInvalidParams", nxcore.ErrInvalidParams)
	nexus.Set("ErrMethodNotFound", nxcore.ErrMethodNotFound)
	nexus.Set("ErrTtlExpired", nxcore.ErrTtlExpired)
	nexus.Set("ErrPermissionDenied", nxcore.ErrPermissionDenied)
	nexus.Set("ErrConnClosed", nxcore.ErrConnClosed)
	nexus.Set("ErrLockNotOwned", nxcore.ErrLockNotOwned)
	nexus.Set("ErrUserExists", nxcore.ErrUserExists)
	nexus.Set("ErrInvalidUser", nxcore.ErrInvalidUser)
	nexus.Set("ErrInvalidPipe", nxcore.ErrInvalidPipe)
	nexus.Set("ErrInvalidTask", nxcore.ErrInvalidTask)
	nexus.Set("ErrCancel", nxcore.ErrCancel)
	nexus.Set("ErrTimeout", nxcore.ErrTimeout)

	println("Nexus Client loaded")
}

func getListOpts(o *js.Object) *nxcore.ListOpts {
	if o == nil {
		o = js.Global.Get("Object").New()
	}
	return &nxcore.ListOpts{
		LimitByDepth: jsObjectGetBool(o, "LimitByDepth", false),
		Depth:        jsObjectGetInt(o, "Depth", -1),
		Filter:       jsObjectGetString(o, "Filter", ""),
	}
}

func getCountOpts(o *js.Object) *nxcore.CountOpts {
	if o == nil {
		o = js.Global.Get("Object").New()
	}
	return &nxcore.CountOpts{
		Subprefixes: jsObjectGetBool(o, "Subprefixes", false),
		Filter:      jsObjectGetString(o, "Filter", ""),
	}
}

func jsObjectGetBool(o *js.Object, k string, d bool) bool {
	if v := o.Get(k); v == js.Undefined {
		return d
	} else {
		return v.Bool()
	}
}

func jsObjectGetInt(o *js.Object, k string, d int) int {
	if v := o.Get(k); v == js.Undefined {
		return d
	} else {
		return v.Int()
	}
}

func jsObjectGetString(o *js.Object, k string, d string) string {
	if v := o.Get(k); v == js.Undefined {
		return d
	} else {
		return v.String()
	}
}
