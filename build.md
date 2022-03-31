# rclone队列自动上传

rclone的单任务上传，适用于出口小宽带，大入口宽带，每次只会有一个任务上传

程序分为客户端和服务端，服务端常运行在后台。每次任务提交，运行客户端命令行提交

可用于其他程序联动，一般使用于qbittorrent、aria2


打包全平台命令
```
goreleaser --snapshot --skip-publish --rm-dist
```