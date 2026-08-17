# BUG-004 复现记录

- 任务类型：diagnosis
- 功能域：version-rollback
- 验证测试：$target
- 隐藏验证：$(TestCanonicalDuplicatePreservesOriginalDraft TestTimeWindowRequiresStrictlyPositiveDuration TestBucketFunctionCanProduceEveryBucket TestHistoricalEvaluationRemainsStableAfterRollbackEdit TestCanceledBatchReturnsWithoutPartialDecisions[4-1])

## 触发命令

`ash
go test -buildvcs=false -count=1 -run "^TestRollbackDraftDoesNotMutatePublishedVersion$" ./internal/service
`

## 预期现象

回滚草稿修改后，历史发布版本的条件值应保持不变；当前基线会被共享切片污染。

## 结果记录

- base：目标验证应失败，公开构建应通过。
- candidate：目标验证继续复现失败；排除该目标后的公开回归和构建应通过。