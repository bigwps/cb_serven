



┌──┤ 🐸🐲 ├─────────🕊️
├─◢◤ 🪐 version 1.0  CB扫描器
├─◢◤ 📌  Go 
├─◢◤🌊 🥲 
├─◢◤ 🌀 FFFFFFFF
└──────────────────────────👻

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

（配合-p在确认有某应用漏洞时可以转发到burp上看详细流量信息）

```powershell
.\xc7.exe -t 127.0.0.1 -p tomcat*spring  -x http://127.0.0.1:8080
```

##### 关闭控制台info信息(-o)

（web扫描开启-o会有未知bug🙄）

```powershell
.\xc7.exe -b 192.168.6.178:5432 -p postgre -o
```

##### dnslog配置在config文件
