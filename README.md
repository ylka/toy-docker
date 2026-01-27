# toy-docker

## 介绍
跟着 https://www.lixueduan.com/posts/docker/mydocker/01-mydocker-run/ 一起手写docker

一些变化：
1 urfave/cli 使用 v3


## 遇到的问题

1. fork/exec /proc/self/exe: operation not permitted

运行时添加 sudo，如 sudo ./toy-docker run -it /bin/sh

2. fork/exec /proc/self/exe: no such file or directory

mount /proc 的时候，切断 mount 传播， 把 / 设成 private