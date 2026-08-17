# BUG-002 复现记录

- 任务类型：diagnosis
- 功能域：rule-validation
- 验证测试：$target
- 隐藏验证：$(TestCanonicalDuplicatePreservesOriginalDraft TestTimeWindowRequiresStrictlyPositiveDuration TestBucketFunctionCanProduceEveryBucket TestHistoricalEvaluationRemainsStableAfterRollbackEdit TestCanceledBatchReturnsWithoutPartialDecisions[2-1])

## 触发命令

`ash
go test -buildvcs=false -count=1 -run "^TestValidateRuleRejectsZeroLengthWindow$" ./internal/service
`

## 预期现象

起止时间相同的规则应被拒绝；当前基线会接受零时长窗口。

## 结果记录

- base：目标验证应失败，公开构建应通过。
- candidate：目标验证继续复现失败；排除该目标后的公开回归和构建应通过。