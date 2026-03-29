# Finance Service

> 最后更新：2026-03-29

## 职责

财务结算层，管理会员订阅（VIP 等级体系）、钱包余额、作者收益结算、提现管理。不管理支付网关对接（由外部支付服务负责）。

## 数据流

```
API Handler → MembershipServiceImpl → WalletRepository → MongoDB
                ↓                        ↓
         AuthorRevenueServiceImpl   AuthorRevenueRepository
              ↓
         ensureWalletCanPay → applyWalletMembershipCharge（钱包扣款）
```

## 约定 & 陷阱

- **钱包扣款原子性**：`ensureWalletCanPay` + `applyWalletMembershipCharge` 必须在同一个事务中，否则会出现余额不一致
- **会员等级映射**：`getLevelFromType` 将会员类型映射为 VIP 等级，新增类型必须同步更新
- **作者收益结算**：收益记录通过 `CalculateEarning` 计算，需注意精度问题（建议用分而非元）
- **提现流程**：`CreateWithdrawalRequest` 创建提现请求，需要审核后才会实际打款
- **会员卡激活**：`ActivateCard` 激活会员卡，同一张卡不能重复激活
