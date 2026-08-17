# BENZHI 验证说明：go-feature-rollout BUG-003

## 项目说明

项目：zhangkui/go-feature-rollout。这是一个 Go 1.22 内存功能开关与灰度发布 Web API。容器使用 golang:1.22，无前端工具链。

## 标准构建、运行和测试命令

`ash
cd '/app' && GOTOOLCHAIN=local go build ./...
cd '/app' && GOTOOLCHAIN=local go run ./cmd
cd '/app' && GOTOOLCHAIN=local go test ./...
`

## Docker 构建和进入容器

`ash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-feature-rollout-bug3-candidate linux/amd64
./build_benzhi_docker.sh go-feature-rollout-bug3-candidate linux/arm64
docker run --rm --platform linux/amd64 -it go-feature-rollout-bug3-candidate bash
docker run --rm --platform linux/arm64 -it go-feature-rollout-bug3-candidate bash
`

## 题目验证命令

任务类型：$type；功能域：$(flag-management rule-validation stable-percentage-rollout version-rollback batch-evaluation[3-1])。

`ash
go test -buildvcs=false -count=1 -run "^TestEvaluatePercentageUsesHundredBuckets$" ./internal/service
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

对于 diagnosis，目标命令预期复现失败，前缀 ! 将该复现转换为验证命令成功；其余构建与回归命令必须通过。

## Bug 复现

详见 [BUG_REPRO.md](BUG_REPRO.md)。