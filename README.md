## Quick start

### 修改环境变量
依赖存在私有仓库，先配置环境变量
```shell
export GH_ACCESS_TOKEN=$git_user:$git_token
```
example:
```shell
export GH_ACCESS_TOKEN=zhangsan:ghp_xxxx
```

> 在[这里](https://github.com/settings/tokens)获取token，获取token的时候选择的scope为repo, workflow


### 准备 `postgres` 数据库

#### 测试环境数据库（推荐）

可以使用测试环境 `pg` 实例，
每个开发者新建一个自己名字命名的 `schema`，
地址可以询问运维或者使用 k8s 命令查看

#### 本地数据库

##### 1. 使用docker安装 postgres

> 建议复用starter-guide下的postgres容器

```sh
## 手动起一个
docker-compose up -d postgres
```


#### 创建表

使用go-migrate工具来维护db schema version


```bash
# 初始
make migrate-create
# 增量更新
make migrate-up

```

### 复制 `.env` 文件

```sh
cp .env.example .env
```

### 使用 `docker` 本地开发

```sh
make run-docker
```

查看日志

```sh
docker-compose logs -f
```

进入容器执行操作

```sh
# 进入容器
docker-compose exec app sh

# 在容器内执行 `protoc` 生成
make gen-proto

# 在容器内执行 `mock` 生成
make mock
```

也可以使用 `vscode` [容器内开发](https://code.visualstudio.com/docs/remote/containers)

## 生成文档

```sh
# 生成 proto 文件时会生成 html 文档
make gen-proto

# 打开 html 文档
open ./doc/index.html
```

## APM

配置 `.env` 文件配置

```diff
- ELASTIC_APM_ENVIRONMENT=development
+ ELASTIC_APM_ENVIRONMENT=<yourname>-local
```

监控 `tidb` sql 查询

```go
// 需要配置 ctx

r.getDB().WithContext(ctx).Find(...)
```


## Project structure

### `cmd/app/main.go`

启动 `internal/app/app.go`.

### `config`

读取 `config.yml` 和 env 环境文件的环境变量，`.env` 中放一些运维需要修改的和敏感的信息，需要export


# 注意
1.rpc
2.代码review
