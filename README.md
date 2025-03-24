# monitor_ssl
监控阿里云，腾讯云，三级域名的https证书有效期



## 配置文件设置

配置文件读取当前路径下的：config.json文件

**concurrence**：多少并发检查https有效期

**alarm_days**：证书不足多少天报警

**email.host**：smtp服务器地址

**email.port**：smtp服务器端口

**email.user**：邮箱登录用户

**email.pass**：邮箱密码

**email.to**：邮件接收者，可多个报警接收人

**aliyun**：阿里云key，需要给与云解析权限

**tencent**：腾讯云key，需要给与云解析权限

**other**：其他需要监控的域名

```json
{
    "concurrence": 100,

    "alarm_days": 30,

    "email":{
        "host": "smtp.163.com",
        "port": "465",
        "user": "1xxxxxxx@163.com",
        "pass": "HMxxxxxxxxxxU",
        "to": ["123456789@qq.com"]
    },
    
    "aliyun": [
        {"key":"阿里云key","secret":"阿里云secret","note":"阿里云账户名"},
        {"key":"LTAI5xxxxxxxxxxxxx","secret":"sglkxxxxxxxx","note":"阿里云账户名"},
    ],

    "tencent": [
        {"SecretId":"腾旭云secret id","SecretKey":"腾讯云secret key","note":"腾讯云账户名"},
        {"SecretId":"AKIxxxxxxxxxxxx","SecretKey":"Woxxxxxxxx","note":"腾讯云账户名"},
    ],

    "other" : [
        "www.baidu.com",
        "www.qq.com"
    ]
}
```

