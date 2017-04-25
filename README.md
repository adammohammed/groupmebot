# GroupMeBot Skeleton
This is a GroupMe Bot skeleton that does the basics for you. The bot is configured using the developer API
information. Once that file is created, the bot is able to start and respond to incoming messages, assuming
your ports are open and your bot is registered to your IP address. To add functionality, simply add functions 
that generate the appropriate repsonses and create a hook with whatever trigger you'd like!
## Setup
### Loads bot from JSON file
1. Create a bot on the [GroupMe Website][1].
2. Make sure and take not of your _bot\_id_ and _group\_id_. 
3. Create a _mybot\_cfg.json_ file.
    ```JSON
    {
      "bot_id": "ffacaer1412....",
      "group_id": "482719...",
      "host": "0.0.0.0",
      "port": "8080"
    }
    ```
Make sure that the bot you created on the [GroupMe Website][1] is registerd to your External IP
and port. Also make sure that whichever port you use, is open on your host.
### Creating your first bot
1. To create the bot first make your folders.
    ```bash
    mkdir $GOPATH/src/github.com/user/mybot
    cd $GOPATH/src/github.com/user/mybot
    ```
2. Create your mybot_cfg.json file shown above with your credentials
       inserted for that *bot_id* and *group_id*
3. Copy/Create a main file similar to the one in this repositories example folder
4. Run this command while in the directory where your main file is located
    ```bash
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

[1]: http://dev.groupme.com/ "Developer GroupMe Website"
