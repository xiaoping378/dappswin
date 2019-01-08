# dappswin

<!-- TODO -->
- 本项目提供链上Restful服务

    由[dep](https://github.com/golang/dep)管理，源码启动方式如下

- 开发启动流程

    - 修改`dappswin.toml`里mysql的配置（用户名和密码）
    - 手动创建database
    
    ```bash
    mysql -u root -pXiaoping@123456 -e "create database if not exists dappswin;"
    ```
    - 启动

    ```bash
    # 下载依赖
    
    make deps

    # 以热编译方式启动
    
    make run
    ```

- glog的高级用法

    ```bash
    # 以V(8)级别打印所有game开头的go文件日志，还打印ico.go的日志

    build/dappswin -vmodule="game*=8,ico=8"

    # 更详细的使用可以 build/dappswin -h 查看

    ```bash

- 同步服务器nginx配置

rsync -avz --delete lottery:/etc/nginx/ ./scripts/nginx


1.项目介绍：
加密乐园（CryptoGaming）是一个全球领先、分散的加密竞猜游戏平台，通过公平、安全、共赢的平台帮助全球加密玩家赢得加密资产，并提供一个成为百万加密鲸鱼的机会。

2.项目优势：
（1）公平、安全；
（2）有趣社交；
（3）多筹码下注；
（4）真人版、VR游戏；
（5）区块链竞猜完整闭环生态价值；
（6）更高赔率、分红、无限彩金；
（7）游戏产生利益公平分配

3.第一款游戏：
加密鲸鱼，一分钟开一次奖，基于EOS区块末尾随机数产生5位数字开奖号码，如果不是数字顺延。
具有9大竞争优势：
（1）最高奖金100万倍；
（2）开奖号码基于EOS区块随机产生，可验证、透明；
（3）即时支付，返还奖金；
（4）保护玩家隐私、安全；
（5）投注即挖矿；
（6）质押分红；
（7）利润竞拍；
（8）鲸鱼榜奖励；
（9）幸运抽奖

4.项目官网、社区信息：* 
官网：dappswin.io* 
中文白皮书：https://dappswin.io/whitepaper/cryptogame-cn-1.0.pdf 
英文白皮书：https://dappswin.io/whitepaper/cryptogame-en-1.0.pdf 
Discord：https://discordapp.com/invite/QWeFf2q
telegram中文：https://t.me/cryptogaming_CN
telegram英文：https://t.me/cryptogaming_EN
Reddit：https://www.reddit.com/user/cryptogamingoffice
Twitter：https://twitter.com/cryptogamingcgg

5.CCN报道：http://t.cn/EGQadnm6.
bitcointalk：https://bitcointalk.org/index.php?topic=5094389.msg49109024#msg49109024

