# BUG-005 复现记录

- 任务类型：bugfix
- 功能域：batch-evaluation
- 验证测试：$target
- 隐藏验证：$(TestCanonicalDuplicatePreservesOriginalDraft TestTimeWindowRequiresStrictlyPositiveDuration TestBucketFunctionCanProduceEveryBucket TestHistoricalEvaluationRemainsStableAfterRollbackEdit TestCanceledBatchReturnsWithoutPartialDecisions[5-1])

## 触发命令

`ash
go test -buildvcs=false -count=1 -run "^TestBatchEvaluateHonorsPreCanceledContext$" ./internal/service
`

## 预期现象

调用前已取消的 context 不应产生任何批量结果；当前基线会返回第一条部分结果和取消错误。

## 结果记录

- base：目标验证应失败，公开构建应通过。
- candidate：目标验证、公开回归和构建应通过。