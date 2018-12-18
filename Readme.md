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
    #下载依赖
    make deps

    #以热编译方式启动
    make run
    ```

