# email_transmission
Use email to transfer files in batches

## 使用方法
* 编译
    > go build -o trans.exe

* 运行
    > ./trans.exe --path '&lt;file path&gt;' --config '&lt;smtp_config.toml&gt;' --sendTo '&lt;email address&gt;' --prefix 'prefix-' --sizeLimit &lt;size&gt;MiB

## 配置文件
```toml
name = "name"
address = "sender@gmail.com"
replyTo = "receiver@gmail.com"
host = "smtp.xxx.com"
port = 465
user = "sender"
password = "password"
keepalive = 60
encryption = true
```