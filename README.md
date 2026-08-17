# go-feature-rollout

## 项目说明

go-feature-rollout 是一个使用 Go 标准库实现的功能开关与灰度发布 Web API。服务支持开关管理、规则配置、稳定百分比分流、不可变版本发布、历史版本回滚和批量决策解释，所有数据保存在内存中。

## 标准命令

```bash
go build ./...
go test ./...
go vet ./...
go run ./cmd
```

## 使用方式

启动服务后监听 `:8080`：

- `POST /flags`：创建功能开关。
- `GET /flags/{key}`：查询功能开关。
- `DELETE /flags/{key}`：归档功能开关。
- `POST /flags/{key}/rules`：添加或更新草稿规则。
- `POST /flags/{key}/publish`：发布草稿版本。
- `POST /flags/{key}/rollback/{version}`：从历史版本创建新草稿。
- `POST /flags/{key}/evaluate`：计算单个主体的开关结果。
- `POST /flags/{key}/batch-evaluate`：批量计算开关结果。

请求和响应均使用 JSON。时间字段使用 RFC3339 格式，规则时间窗口采用起点包含、终点不包含的语义。
