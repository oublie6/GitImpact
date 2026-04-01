# 开发说明

## 后端开发
- 技术栈：Go + Gin + GORM。
- 默认采用 vendor 模式：`GOFLAGS=-mod=vendor`。

常用命令：
- `./scripts/init-dev.sh`
- `make build`
- `make test`
- `make build-linux-amd64`
- `docker build -t gitimpact/backend:test .`

## Windows PowerShell 交叉编译 Linux x64
在仓库根目录执行：

```powershell
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:GOFLAGS="-mod=vendor"
cd backend
go build -trimpath -ldflags "-s -w" -o ../bin/gitimpact-backend-linux-amd64 ./cmd/server
```

也可以直接执行：

```powershell
make build-linux-amd64
```

## Docker 构建验证
在仓库根目录执行：

```bash
docker build -t gitimpact/backend:test .
```

该构建会使用根目录 `Dockerfile`，并在容器内基于 `backend/vendor` 进行 `GOFLAGS=-mod=vendor` 构建，不联网拉取 Go 依赖。
