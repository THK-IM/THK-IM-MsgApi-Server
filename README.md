# THK-IM-MsgApi-Server

## 启动服务

- msg api服务

```
go run main.go --config-path etc/msg_api_server.yaml

```

## 构建镜像

```
docker build -t thk-im/msg-api-server:v1  -f ./deploy/.Dockerfile .
```

