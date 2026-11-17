# freedom-list — P78-RELEASE-ZIP-BINDING

切片：`RELEASE-ZIP-BINDING`（P78，v0.1.0 Release 创建 + 最终 ZIP 构建 + 发行记录绑定 + 安装生命周期复验，实现类 `[SLICE]`，M-L 档）。本清单由主控在 ⑥ 之后、⑦ 之前定稿（§1.7.3）；步骤 0 先落初稿。

## 一、已行使的执行自由度（初稿，待 ⑦ 前定稿）

| 编号 | 实质性选择 | 依据 |
|---|---|---|
| F-1 | 方案书 `design-r1.md` 单语（中文），不做双语；双语报告留 ②e 阶段 | 起跑消息 §三（P77 / P76 先例） |
| F-2 | 证据包目录沿用 P77 体例（`01-arch` / `02-batch` / `04-reviews` / `06-gates`），复验证据落 `05-verify` 子目录 | 起跑消息 §四 ②e |
| F-3 | 构建路径（CI 触发 vs 本机同参数）由 ① 判定，主控只执行判定结果；若两者不可满足则停下升级 | 起跑消息 §三 Q1 + §十一 第 2 条 |
| F-4 | 最终 ZIP 的包装层（是否保留 `prfrail/` 顶层目录）严格照 `INSTALLATION{,_EN}` 已冻结的「单顶层目录 + 四个发行文件」口径 | ADR-009 + INSTALLATION §2 |
| F-5 | Release notes 正文由主控起草并落 `tmp/` 后上传，不新增仓库内 ADR / OB | 起跑消息 §七 边界 |
| F-6 | 第二层文档按 ① Q5 实读结果决定改写、加法式指针或不动；升降级在方案书如实披露 | 起跑消息 §八末尾 |
| F-7 | 批量辅助两次的产出只作事实枚举输入，不进方案书正文的判定栏位 | 起跑消息 §六 批量辅助纪律 |

## 二、未行使的自由度（明确不做）

- 不改产品代码（`.go` / `go.mod`）；构建失败按 §1.7.2 停下报实况。
- 不改 `ci.yml`；workflow_dispatch 触发仅为外部动作。
- 不改 ADR-009 ZIP 结构、ADR-006 发布信任；不做代码签名。
- 不改 S1 Exit 八门判据；`p57-closeout-assessment{,_EN}` 门表与 `s1-exit{,_EN}` 门表零改动。
- 不改 `S1_SUPPORT_MATRIX` 平台清单与支持层级（绑定为加法式）；不做 Linux 正式发布。
- 不改 ADR-010 EOL 通知方式；`ADR_REGISTER:62`「当前仍不是已发布 SLA」默认不翻。
- 不新增 ADR / OB（除非用户裁定）；不改准则正文（维持 v1.69）；不改 v50–v61 任何快照。
- 不改 `CONTRACTS` / `SECURITY` / `ARCHITECTURE` / `TEST_STRATEGY` / `BUSINESS_WORKFLOWS` / `CODING_CONVENTIONS` / `DOCUMENTATION_PLAN` / `RFC-*` / schema / fixtures。
- 不做任何 gitee 动作；Release 发布对象仅 GitHub origin。

## 三、定稿规则

⑥ 结论与 ⑦ 输入包范围确定后，在本文件回填「⑥ 之后的最终选择」并重算 `SHA256SUMS.txt`；⑧b 再把同一清单回写进报告（二者内容必须一致）。

## 四、定稿（⑥ 之后、⑦ 之前）

- （⑥ 之后回填）
