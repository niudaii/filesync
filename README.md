# filesync
#### 简介

server 和 worker 文件同步工具。

#### 使用

```
文件同步工具 by zp857

Usage:
  filesync [command]

Available Commands:
  help        Help about any command
  server      文件同步服务端
  worker      文件同步客户端

Flags:
      --auth string   auth (default "zp857")
      --dir string    dir to file sync (default "./")
  -h, --help          help for filesync
      --host string   host
  -p, --port string   port (default "5001")

Use "filesync [command] --help" for more information about a command.
```

```
./filesync server --host 0.0.0.0 --dir resource --black webscan --monitor
```

![image-20221230102654977](https://nnotes.oss-cn-hangzhou.aliyuncs.com/notes/image-20221230102654977.png)

```
./filesync worker --host xxx.xxx.xxx.xxx --dir resource --timer 3600
```

![image-20221230102724929](https://nnotes.oss-cn-hangzhou.aliyuncs.com/notes/image-20221230102724929.png)
