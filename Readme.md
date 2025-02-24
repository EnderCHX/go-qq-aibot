# QQ机器人接入Deepseek Ai

## QQ机器人使用OneBot11协议，ws正向连接

## 根目录创建config.json文件，内容如下：

```json
{
    "qqbot" : {
        "ws_addr" : "ws://ip:port",
        "key" : "token",
        "nickname" : ["喵喵"],
        "command_prefix" : "#",
        "superusers" : [123456789,987654321]
    },
    "deepseek" : {
        "api_url" : "https://api.deepseek.com/chat/completions",
        "api_key" : "你的api key",
        "model" : "deepseek-reasoner",
        "sys_prompt" : "你是一只可爱猫娘"
    },
    "websearch" : {
        "api_url" : "your searxng api url",
    }
}
```
