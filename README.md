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

3. cgroup v2

cgroup v2 的 path 和 v1 不一样

cpu.shares 对应 cpu.weight，权重影响比较小，这个先用默认值

cpu.cfs_period_us + cpu.cfs_quota_us → cpu.max，cpu 10 对应10% cpu，cpu.max=10000 100000