

本项目为go语言练习项目，兼容xray poc ，但没有写search output规则（懒），代码写了详细的注释，欢迎大家指正。



```
┌──┤ 🐸🐲 ├─────────🕊️
├─◢◤ 🪐 version 1.0  CB扫描器
├─◢◤ 📌  Go 
├─◢◤🌊 🥲 
├─◢◤ 🌀 FFFFFFFF
└──────────────────────────👻
```



#### 1.WEB扫描

```powershell
.\xc7.exe -t 127.0.0.1                
.\xc7.exe -t 127.0.0.1 -p tomcat*spring  #只扫描tomcat和spring漏洞
```

#### 2.数据库弱口令爆破

（需数据库支持远程连接）

```powershell
.\xc7.exe -b 192.168.6.178:5432 -p postgre
```

#### 3.其他

##### web扫描从文件读取目标（-f）

```powershell
.\xc7.exe -f ./11.txt -p tomcat*spring 
```

##### web扫描代理(-x)

只写了http代理

（配合-p在确认有某应用漏洞时可以转发到burp上看详细流量信息）

```powershell
.\xc7.exe -t 127.0.0.1 -p tomcat*spring  -x http://127.0.0.1:8080
```

##### 关闭控制台info信息(-o)

（web扫描开启-o会有未知bug🙄）

```powershell
.\xc7.exe -b 192.168.6.178:5432 -p postgre -o
```

##### dnslog配置在config文件（必须要配置）



#### 参考-感谢

[chaitin/xray: 一款完善的安全评估工具，支持常见 web 安全问题扫描和自定义 poc | 使用之前务必先阅读文档 (github.com)](https://github.com/chaitin/xray)

[shadow1ng/fscan: 一款内网综合扫描工具，方便一键自动化、全方位漏扫扫描。 (github.com)](https://github.com/shadow1ng/fscan)

[jjf012/gopoc: 用cel-go重现了长亭xray的poc检测功能的轮子 (github.com)](https://github.com/jjf012/gopoc)

