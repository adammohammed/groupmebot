# GroupMeBot Framework
This is a GroupMe Bot Framework that does the basics for you. The bot is configured using the developer API
information. Once that file is created, the bot is able to start and respond to incoming messages, assuming
your ports are open and your bot is registered to your IP address. To add functionality, simply add functions 
that generate the appropriate repsonses and create a hook with whatever trigger you'd like!
## Setup
### Loads bot from JSON file
1. Create a bot on the [GroupMe Website][1].
2. Make sure and take note of your _bot\_id_ and _group\_id_. 
3. Create a _mybot\_cfg.json_ file.
```javascript
{
  "bot_id": "your_bot_id",
  "group_id": "your_group_id",
  "host": "0.0.0.0",
  "port": "8080"
}
```
Make sure that the bot you created on the [GroupMe Website][1] is registerd to your External IP
and port. Also make sure that whichever port you use, is open on your host.
### Installing Go
This project uses the programming language [Go](https://golang.org).
Make sure you install it and set your GOPATH environment variable to the location of the
directory that you are working.

### Creating your first bot
1. To create the bot first make your folders.
```sh
mkdir $GOPATH/src/github.com/user/mybot
cd $GOPATH/src/github.com/user/mybot
```
2. Create your mybot_cfg.json file shown above with your credentials
       inserted for that *bot_id* and *group_id*
3. Copy/Create a main file similar to the one in this repositories example folder
4. Run this command while in the directory where your main file is located
```sh
go get -u -v github.com/adammohammed/groupmebot
go install
```
5. Finally you can run your bot using the `go run main.go` command
### Creating plugins
Defining hooks is simple. The hooks signature needs to be as follows:
```go
func my_hook(msg groupmebot.InboundMessage) (string) {
        resp := fmt.Sprintf("Hi there, %v.", msg.Name)
        return resp
}
```

This function must accept only an Inbound message as input.
The output must be a string. The actual body of the fuction can do whatever you please
but the signature is vital to be able to add it to the list of hooks.

Adding items to the hooks is done as shown below. Assume that these are function names for functions
defined with the signature we defined earlier.
```go
bot.AddHook("Hi!$", my_hook)
```

The first parameter to the AddHook method is the regular expression trigger that the message text must match.
In this example, if the message text ends with "Hi!", the bot will send a response.
If the bot finds the incoming message matches the trigger expression, the function is executed. When a 
message is received by the bot, it checks to see if the text matches any available hooks added by the 
AddHooks method.

## Future changes
- Make it possible to dynamically populate hooks from plugin directory.
- Allow hooks to be removed/updated

## License ##
Copyright (c) 2017 Adam Mohammed

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

[1]: http://dev.groupme.com/ "Developer GroupMe Website"
