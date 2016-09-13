Nexus Javascript / Node.js Client
=================================

Nexus client for web browser/node.js built using gopherjs, wrapping around the golang nexus client

# Requirements
  * [GopherJS](https://github.com/gopherjs/gopherjs)
    * ```go get github.com/gopherjs/gopherjs```
  
## Requirements for node.js:
  * WebSocket module ['ws'](https://github.com/websockets/ws)
    * ```npm -g install ws```

# Build
```bash
$ gopherjs build
```

# API
// WIP: Missing a proper documentation. These are pseudo-go headers

```javascript
    // On browsers
    nexus.dial(url, callback) (NexusConnection, error)

    // On Node.js
    var nexus = require("./nxjs.js")
    nexus.dial(url, callback) (NexusConnection, error)

    NexusConnection Object:
        func login(user string, pass string, callback)
        
        func taskPush(method string, params interface{}, timeout int, callback)
        func taskPull(prefix string, timeout int, callback) (Task, error)
        func taskList(prefix string, limit int, skip int, callback)
        
        func userCreate(user string, pass string, callback)
        func userDelete(user string, callback)
        func userDelTags(user string, prefix string, tags []string, callback)
        func userSetPass(user string, pass string, callback)
        func userSetTags(user string, prefix string, tags map[string]{}, callback)
        func userAddTemplate(user string, template string, callback)
        func userDelTemplate(user string, template string, callback)
        func userAddWhitelist(user string, ip string, callback)
        func userDelWhitelist(user string, ip string, callback)
        func userAddBlacklist(user string, ip string, callback)
        func userDelBlacklist(user string, ip string, callback)
        func userSetMaxSessions(user string, max int, callback)
        func sessionKick(connId string, callback)
        func sessionReload(connId string, callback)
        func userList(prefix string, limit int, skip int, callback)

        func nodeList(limit int, skip int, callback)
        func node(callback)

        func pipeCreate(opts interface{}, callback) (Pipe, error)
        func pipeOpen(id string, callback) (Pipe, error)

        func topicPublish(channel string, msg interface{}, callback)
        func topicSubscribe(pipe Pipe, channel string, callback)
        func topicUnsubscribe(pipe Pipe, channel string, callback)

        func lock(lock string, callback)
        func unlock(lock string, callback)

        func exec(method string, params interface{}, callback)
        func execNoWait(method string, params interface{}, callback)

        func cancel(callback)
        func closed(callback)
        func version(callback)
        func nexusVersion(callback)
        func ping(timeout int, callback)
        
    Task Object:
        func sendResult(res interface{}, callback)
        func sendError(code int, msg string, data interface{}, callback)
        field Path
        field Method
        field Params
        field Tags

    Pipe Object:
        func close(callback)
        func read(max int, timeout int, callback)
        func write(msg {}, callback)
        func id(callback) string

```

Functions can receive zero, one or two parameters at the end for callbacks.

If there is only one callback parameter, it should be a function with two arguments:
```javascript
  taskPull("prefix", 60, function(result, error) {
    console.log("This is the result:", result)
    console.log("Error received:", error)
  }
```

With two callback parameters, one will receive the result and the other the error:
```javascript
  login("user", "pass", function(result) {
    console.log("Logged in!")
  }, function(error) {
    console.log("Couldn't login!:", error)
  })
```


# Examples

## Pull a task from a browser
```javascript
// The module will set dial as a global function when loaded
nexus.dial("wss://localhost.n4m.zone", function(nc, err){

  // Login to nexus
  nc.login("dummyUser", "dummyPassword", function(){
  
    // Success! Now pull a task
    nc.taskPull("test.prefix", 5, function(task, err){
    
      // Great! Just return an OK
      console.log(task, err);
      task.SendResult("OK");
    })
    
  })
})
```

## Subscribe a pipe to a channel

```javascript
var nexus = require("./nxjs.js")

nexus.dial("wss://localhost.n4m.zone", function(nc, err){
  nc.login("dummyUser", "dummyPass", function(){
  
    // Create a pipe
    nc.pipeCreate({"len": 100}, function(pipe, e){
    
      // Subscribe the pipe to the channel
      nc.chanSubscribe(pipe, "temperatures",
      
        //Subscription succeeded
        function(){
          console.log("Subscribed pipe", pipe.id(), "to channel temperatures")
        
          pipe.read(10, 60, function(msgs, err){
            console.log("Received messages:", msgs)
          })
        },
        
        // Subscription failed
        function(err){ console.log("Error subscribing the pipe to the channel:", err)})
    })
  })
})
```