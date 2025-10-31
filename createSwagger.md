# Swagger API 文档集成

## 概述

本项目已集成 Swaggo，可以自动生成 OpenAPI/Swagger 文档。

## 安装和配置

### 1. 安装 swag 命令行工具

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. 安装项目依赖

```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

## 使用方法

### 1. 生成 Swagger 文档

```bash
# 使用 Makefile
make swagger

# 或直接使用 swag 命令
swag init -g main.go --output docs
```

### 2. 启动项目

```bash
go run main.go
```

### 3. 访问 Swagger UI

启动项目后，在浏览器中访问：

- **Swagger UI**: http://localhost:8082/swagger/index.html
- **API 文档 JSON**: http://localhost:8082/swagger/doc.json

## API 文档注解说明

### 通用注解格式

```go
// HandlerName 处理器名称
// @Summary 简短描述
// @Description 详细描述
// @Tags 标签分组
// @Accept 接受的内容类型
// @Produce 返回的内容类型
// @Param 参数说明
// @Success 成功响应
// @Failure 失败响应
// @Router 路由路径 [HTTP方法]
func HandlerName(c *gin.Context) {
    // 处理器实现
}
```

### 参数类型说明

- `query` - URL 查询参数
- `path` - URL 路径参数
- `header` - HTTP 头参数
- `body` - 请求体参数
- `formData` - 表单参数

### 示例

```go
// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户ID获取详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param token header string true "认证令牌"
// @Success 200 {object} User "用户信息"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Router /api/users/{id} [get]
func GetUserInfo(c *gin.Context) {
    // 实现代码
}
```

## 当前已集成的 API

### 1. IP 信息 (`/ip`)

- **GET** `/ip` - 获取 IP 地址信息

### 2. 邮箱验证 (`/api/send-email`, `/api/verify-code`)

- **POST** `/api/send-email` - 发送邮箱验证码
- **POST** `/api/verify-code` - 验证邮箱验证码

### 3. 图形验证码 (`/api/captcha`, `/api/verify-captcha`)

- **POST** `/api/captcha` - 获取图形验证码
- **POST** `/api/verify-captcha` - 验证图形验证码

### 4. B 站 API (`/api/bili-follow-*`)

- **GET** `/api/bili-follow-anime` - 获取追番列表
- **GET** `/api/bili-follow-movie` - 获取追剧列表

### 5. 静态文件 (`/api/static-files`)

- **GET** `/api/static-files` - 获取静态文件列表

## 开发工作流

### 1. 添加新的 API 端点

1. 在对应的 handler 文件中添加 Swagger 注解
2. 重新生成文档：`make swagger`
3. 重启项目查看更新

### 2. 修改 API 文档

1. 更新 handler 中的注解
2. 重新生成文档：`make swagger`
3. 刷新浏览器查看变更

### 3. 自定义响应模型

```go
type UserResponse struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"张三"`
    Email string `json:"email" example:"user@example.com"`
}
```

### 4. 安全配置

项目已配置 JWT 认证：

```go
// @Security ApiKeyAuth
// @Router /api/protected [get]
```

## 常用命令

```bash
# 生成并运行
make swagger-run

# 仅生成文档
make swagger

# 运行项目
make run

# 编译项目
make build

# 清理文件
make clean
```

## 注意事项

1. 每次修改 API 注解后都需要重新生成文档
2. `docs/` 目录中的文件是自动生成的，不要手动修改
3. 注解语法严格，注意格式和拼写
4. 建议在开发时使用 `air` 实现热重载

## 故障排除

### 问题 1: swag 命令不存在

```bash
# 解决方案
go install github.com/swaggo/swag/cmd/swag@latest
```

### 问题 2: 文档生成失败

- 检查注解语法是否正确
- 确认所有引用的结构体都已定义
- 查看命令行错误输出

### 问题 3: Swagger UI 无法访问

- 确认项目已启动
- 检查端口是否正确
- 验证路由是否正确配置
