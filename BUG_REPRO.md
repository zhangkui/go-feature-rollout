# BUG-001 复现记录

- 任务类型：bugfix
- 功能域：flag-management
- 验证测试：$target
- 隐藏验证：$(TestCanonicalDuplicatePreservesOriginalDraft TestTimeWindowRequiresStrictlyPositiveDuration TestBucketFunctionCanProduceEveryBucket TestHistoricalEvaluationRemainsStableAfterRollbackEdit TestCanceledBatchReturnsWithoutPartialDecisions[1-1])

## 触发命令

`ash
go test -buildvcs=false -count=1 -run "^TestCreateFlagRejectsCanonicalDuplicate$" ./internal/service
`

## 预期现象

使用规范化后相同但原始文本不同的开关键创建第二次时应返回冲突，不能覆盖已有定义。

## 结果记录

- base：目标验证应失败，公开构建应通过。
- candidate：目标验证、公开回归和构建应通过。