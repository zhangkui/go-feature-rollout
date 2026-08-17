# BUG-003 复现记录

- 任务类型：bugfix
- 功能域：stable-percentage-rollout
- 验证测试：$target
- 隐藏验证：$(TestCanonicalDuplicatePreservesOriginalDraft TestTimeWindowRequiresStrictlyPositiveDuration TestBucketFunctionCanProduceEveryBucket TestHistoricalEvaluationRemainsStableAfterRollbackEdit TestCanceledBatchReturnsWithoutPartialDecisions[3-1])

## 触发命令

`ash
go test -buildvcs=false -count=1 -run "^TestEvaluatePercentageUsesHundredBuckets$" ./internal/service
`

## 预期现象

99% 放量时桶 99 的主体应关闭；当前基线报告不到桶 99 且会错误启用。

## 结果记录

- base：目标验证应失败，公开构建应通过。
- candidate：目标验证、公开回归和构建应通过。