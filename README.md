# nat-proxy
## 内外服务代理转发工具，通过该工具可以将内外的服务对外发布，即通过公网地址端口代理转发内外的服务，功能类似frp
## server端运行在公网服务器
## client运行在内外机器

## 配置文件详见bin/config.ini说明


## 启动
### 1 启动服务端
./proxy_server config.ini


### 2 启动客户端
./proxy_client2 config.ini