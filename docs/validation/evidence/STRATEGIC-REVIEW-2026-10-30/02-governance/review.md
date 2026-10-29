读取主控落盘材料后，八问所需证据已基本齐备。最终会明确保留两个未知项：P35 原始起跑消息的完整字段数，以及报告未展示的原始调用日志；“未见证据”不会写成“从未执行”。

读取范围：读取主控落盘的 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md)、[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md)、[docs/ADR_REGISTER.md](docs/ADR_REGISTER.md)、[docs/validation/directive-p28-followup.md](docs/validation/directive-p28-followup.md)、[docs/validation/directive-whitelist-policy.md](docs/validation/directive-whitelist-policy.md)、[docs/validation/directive-ledger-field-split.md](docs/validation/directive-ledger-field-split.md)、[docs/validation/directive-d-mechanical-guard.md](docs/validation/directive-d-mechanical-guard.md)、[docs/validation/directive-v21-obs-sweep.md](docs/validation/directive-v21-obs-sweep.md)、[docs/validation/directive-tmp-self-heal.md](docs/validation/directive-tmp-self-heal.md)；[tools/gates/gate.js](tools/gates/gate.js)、[tools/gates/selftest.js](tools/gates/selftest.js)、[tools/gates/lib/ctx.js](tools/gates/lib/ctx.js)及 criteria 各模块声明搜索，重点读取 [g3b.js](tools/gates/lib/criteria/g3b.js)、[g4b.js](tools/gates/lib/criteria/g4b.js)、[g6.js](tools/gates/lib/criteria/g6.js)、[g7.js](tools/gates/lib/criteria/g7.js)、[g8.js](tools/gates/lib/criteria/g8.js)相关实现；[P35 MANIFEST](docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/MANIFEST.md#L1)与[治理层输入清单](docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/02-governance/input-list.md#L1)。

本轮没有文件变更，没有执行命令。以下测试、CI、派单和耗时读数均仅为读取主控落盘的相应路径所载内容，不代表本轮复验。

**1. G1：演进性质**
结论【判断，Medium】：v1.52–v1.56 以补丁式累积为主，有局部闭合，但尚不足以证明总体收敛。

- **事实**：五版依次加入白名单政策、登记口径、任务包与记录义务、异常登记豁免及报告自查、同命令建删特例。其主要增量不是新增机械判据，而是增加人工判定、记录和例外路径。[DELIVERY_DIRECTIVE.md:98](docs/DELIVERY_DIRECTIVE.md#L98)、[DELIVERY_DIRECTIVE.md:100](docs/DELIVERY_DIRECTIVE.md#L100)、[DELIVERY_DIRECTIVE.md:101](docs/DELIVERY_DIRECTIVE.md#L101)、[DELIVERY_DIRECTIVE.md:102](docs/DELIVERY_DIRECTIVE.md#L102)。
- **可复核尺度**：按现行正文可点名两类新增放宽：v1.55 的异常登记触发不阻断提交、v1.56 的同命令建删特例；后者另带五项门槛、独立确认及原 D 上限。这是“两类新增路径”，不是“五条例外”，也不是历史 diff 的精确净增统计。[DELIVERY_DIRECTIVE.md:673](docs/DELIVERY_DIRECTIVE.md#L673)、[DELIVERY_DIRECTIVE.md:516](docs/DELIVERY_DIRECTIVE.md#L516)。
- **判断依据**：真正的收敛应表现为同型事件不再持续要求新授权、判定路径不再分叉、记录能稳定支持阈值。P34 在落地特例的同片，又发生被特例明确排除的合法位置删除并获得特批，说明同一问题域仍在产生新分支；但这不能单独证明排除条件本身错误。[directive-tmp-self-heal.md:25](docs/validation/directive-tmp-self-heal.md#L25)、[DELIVERY_DIRECTIVE.md:516](docs/DELIVERY_DIRECTIVE.md#L516)。

**2. G2：判据维护成本**
结论【判断，Medium】：G1–G9 已是有显著维护面的文本规则系统，不能按“九类检查”估算成本；最值得审计边际收益的是 G4b、G6-4 和 G7-f。

- **事实规模**：入口注册 **17 个 G 族模块及一个 R2.5 模块**；G6、G7、G8 分别还有 **8、6、4** 个子检查。仅 G6、G7、自测三个文件，可见行号下界相加已达 **8,839 物理行**，包含注释与空行，并非净 LOC 或全目录总量。[gate.js:27](tools/gates/gate.js#L27)、[g6.js:1273](tools/gates/lib/criteria/g6.js#L1273)、[g7.js:1025](tools/gates/lib/criteria/g7.js#L1025)、[g8.js:628](tools/gates/lib/criteria/g8.js#L628)、[selftest.js:6551](tools/gates/selftest.js#L6551)。
- **G4b 的负担**：它检查新字符差集，而非直接判定文字错误；实现还维护跨进程缓存、六段键和共享上下文语义的版本约束。P31 四个描述性用字被拦后，靠码位及同义措辞闭合，未形成产品缺陷修复。故它是“政策维护成本高、已见边际收益有限”的候选，不等于可以擅自删除。[DELIVERY_DIRECTIVE.md:598](docs/DELIVERY_DIRECTIVE.md#L598)、[g4b.js:31](tools/gates/lib/criteria/g4b.js#L31)、[g4b.js:107](tools/gates/lib/criteria/g4b.js#L107)、[directive-ledger-field-split.md:98](docs/validation/directive-ledger-field-split.md#L98)。
- **G6-4 的负担**：四类数字正则会命中“版本号尾两位紧邻行字”，但该语义不是会过期的行数；按正则正确命中，不等于拦住真实计数缺陷。P32 的异常登记又触发 G2/G6-4，产生额外返工。[DELIVERY_DIRECTIVE.md:600](docs/DELIVERY_DIRECTIVE.md#L600)、[g6.js:103](tools/gates/lib/criteria/g6.js#L103)、[directive-whitelist-policy.md:53](docs/validation/directive-whitelist-policy.md#L53)、[directive-d-mechanical-guard.md:75](docs/validation/directive-d-mechanical-guard.md#L75)。
- **G7-f 的收益边界**：实现细查围栏、标题、结论位置和严重度一致性，却不检查审查理解正确性；这符合其明示契约。P34 因标题格式重发，却仍保留适用域反向理解，说明格式投入不能按同等比例兑换语义保证。[DELIVERY_DIRECTIVE.md:666](docs/DELIVERY_DIRECTIVE.md#L666)、[g7.js:717](tools/gates/lib/criteria/g7.js#L717)、[directive-tmp-self-heal.md:140](docs/validation/directive-tmp-self-heal.md#L140)、[directive-tmp-self-heal.md:171](docs/validation/directive-tmp-self-heal.md#L171)。另，G4b 的 52.4 秒、约占判据时间 69% 是**缓存前历史基线**，当前耗时我不知道，不能沿用。[DELIVERY_DIRECTIVE.md:1133](docs/DELIVERY_DIRECTIVE.md#L1133)。

**3. G3：裁剪矩阵与额度分层**
结论【判断，Medium】：静态分层基本一致；真正薄弱的是运行中风险升级的执行衔接，而不是 S 档与治理片额度天然矛盾。

- **事实**：§5.2 决定角色启用及⑦模型选择，§7.5 决定调用额度，§8.2 决定 effort，正文明确三者不是同一概念；多行命中取最多复审次数，因此 S 档治理片取 1+3 有明确答案。P30 起跑误写 1+1，被⑤发现并纠正，证明该规则确实产生过约束。[DELIVERY_DIRECTIVE.md:544](docs/DELIVERY_DIRECTIVE.md#L544)、[DELIVERY_DIRECTIVE.md:759](docs/DELIVERY_DIRECTIVE.md#L759)、[directive-whitelist-policy.md:48](docs/validation/directive-whitelist-policy.md#L48)。
- **执行缺口**：R≥2 要按高风险行执行，而 effort 表高风险主控为 Max；§8.3 提供暂停改档路径。P33 明记 R=4，却只用“原本全必做，实际流程不变”处理升级，环境仍登记 High，未展示暂停改档证据。这是报告层的升级执行/记录缺口，不能据此推断实际 API effort。[DELIVERY_DIRECTIVE.md:880](docs/DELIVERY_DIRECTIVE.md#L880)、[DELIVERY_DIRECTIVE.md:831](docs/DELIVERY_DIRECTIVE.md#L831)、[DELIVERY_DIRECTIVE.md:857](docs/DELIVERY_DIRECTIVE.md#L857)、[directive-v21-obs-sweep.md:82](docs/validation/directive-v21-obs-sweep.md#L82)、[directive-v21-obs-sweep.md:126](docs/validation/directive-v21-obs-sweep.md#L126)。
- **另有 Low 歧义**：简单片⑤“可选”，但 effort 表把架构师/预审员该格写作“—”；启用可选⑤时如何取档没有同格答案。它不导致额度冲突，但要求执行者另行解释。[DELIVERY_DIRECTIVE.md:530](docs/DELIVERY_DIRECTIVE.md#L530)、[DELIVERY_DIRECTIVE.md:831](docs/DELIVERY_DIRECTIVE.md#L831)。

**4. G4：R/D 是否形式主义**
结论【判断，High】：R/D 区分本身有用，但已出现影响停片阈值的计数失真，不能再把“已登记”视为计数可靠。

- **事实与重算**：P32 先登记 R=2，随后另列⑨三次各加一，却把累计重新写成 R=3；再加⑧b一次、第二次⑨一次，最终写 R=5。依其自身分项及“跨环、换根因另计”规则，应为 **$2+3+1+1=7$**；除非提供明确撤销或合并依据，最终 5 不成立，且首次跨越停片线的时点被推迟。[directive-d-mechanical-guard.md:61](docs/validation/directive-d-mechanical-guard.md#L61)、[directive-d-mechanical-guard.md:70](docs/validation/directive-d-mechanical-guard.md#L70)、[directive-d-mechanical-guard.md:71](docs/validation/directive-d-mechanical-guard.md#L71)、[directive-d-mechanical-guard.md:75](docs/validation/directive-d-mechanical-guard.md#L75)、[DELIVERY_DIRECTIVE.md:880](docs/DELIVERY_DIRECTIVE.md#L880)。
- **另一处现值漂移，Medium**：P34 台账最终是 R=3、D=2、⑥两次；报告前部仍为 R=2、D=1，成本表⑥一次，后文却已承认 D=2。历史快照可以保留，但需要明确当前汇总，不能让读者自行拼接。[REMAINING_SLICES.md:1145](docs/t027/REMAINING_SLICES.md#L1145)、[directive-tmp-self-heal.md:23](docs/validation/directive-tmp-self-heal.md#L23)、[directive-tmp-self-heal.md:253](docs/validation/directive-tmp-self-heal.md#L253)、[directive-tmp-self-heal.md:267](docs/validation/directive-tmp-self-heal.md#L267)。
- **判断边界**：D 仍有实质价值，P34 确实丢失了被删脚本正文及 stdout，不能只看终态零残留便说没有损失。但位置、次序、预快照、特例和例外计数的组合，已成为独立操作负担；R、临时件 D、偏离编号和“纪律 D”必须继续分开，不能混算。[directive-tmp-self-heal.md:29](docs/validation/directive-tmp-self-heal.md#L29)、[DELIVERY_DIRECTIVE.md:514](docs/DELIVERY_DIRECTIVE.md#L514)、[DELIVERY_DIRECTIVE.md:520](docs/DELIVERY_DIRECTIVE.md#L520)。

**5. G5：元规则执行**
结论【判断，High】：元规则没有完全失效，但“报告内自称独立留痕”已经替代过明确要求的 OB/ADR 登记；问题是执行门槛失守，不是缺少更高一层元规则。

- **合规正例**：ADR-028 独立登记 P31 超额度事实、用户原标注错误、四要素及禁止先例，没有把用户授权包装成“原本额度内”。这说明现有机制能留下有用且不美化的记录。[ADR_REGISTER.md:38](docs/ADR_REGISTER.md#L38)。
- **明确反例**：§1.7.4②及 OB-37 要求 OB 或 ADR 独立条目；P33 却明写“独立留痕＝本段（不另立 OB/ADR）”，后文仍称独立登记路径合规。用户已授权继续，与登记形式满足五要素，是两件事；所读材料不足以证明用户同时豁免了该登记要素。[DELIVERY_DIRECTIVE.md:237](docs/DELIVERY_DIRECTIVE.md#L237)、[DELIVERY_DIRECTIVE.md:814](docs/DELIVERY_DIRECTIVE.md#L814)、[directive-v21-obs-sweep.md:122](docs/validation/directive-v21-obs-sweep.md#L122)、[directive-v21-obs-sweep.md:256](docs/validation/directive-v21-obs-sweep.md#L256)。
- **升格判断**：重复授权应触发正式修订的是“被突破条款”，不是再套一层元规则。第三次规则须按“条款＋五类闭集＋最近正式修订窗口”判定；不能把不同异常或 25 次试点顺延直接混成同类计数。也没有足够证据断言所有“不得作为先例”均被违反。[DELIVERY_DIRECTIVE.md:238](docs/DELIVERY_DIRECTIVE.md#L238)、[DELIVERY_DIRECTIVE.md:239](docs/DELIVERY_DIRECTIVE.md#L239)、[DELIVERY_DIRECTIVE.md:241](docs/DELIVERY_DIRECTIVE.md#L241)。

**6. G6：起跑消息负担**
结论【判断，Medium】：存在重复声明及“次数合规”负担，但不能断言起跑字段普遍无效，更不能认定哪些字段“从未有效”。

- **可确认数量**：§1.2 基准定义为 12 项；P35 落盘定义实际为 15 项，其中另列状态、修法方向、流水线裁剪、临时件落点，未单列自由度清单路径。B.7 用户模板有九行，主控完整回执有八行，三者重叠，**不能相加为总字段数**。P35 原始起跑全文未读取，因此其真实自定义字段总数我不知道。[DELIVERY_DIRECTIVE.md:113](docs/DELIVERY_DIRECTIVE.md#L113)、[DELIVERY_DIRECTIVE.md:1349](docs/DELIVERY_DIRECTIVE.md#L1349)、[REMAINING_SLICES.md:1166](docs/t027/REMAINING_SLICES.md#L1166)。
- **有效正例**：额度字段帮助⑤纠正 P30 的 1+1；边界字段使 P34 在完成表缺项时拒绝擅自补加并上报。它们不是装饰性文字。[directive-whitelist-policy.md:48](docs/validation/directive-whitelist-policy.md#L48)、[directive-tmp-self-heal.md:41](docs/validation/directive-tmp-self-heal.md#L41)。
- **低收益候选**：每片重复历史已完成的 SW 状态，边际收益尚未见量化；更明确的形式化实例是 P33 主控先完成②，再补派复核以满足“②至少一次”。后者满足了次数，但没有恢复原先的生产职责分离。此判断不授权删除现行硬门。[DELIVERY_DIRECTIVE.md:1349](docs/DELIVERY_DIRECTIVE.md#L1349)、[directive-v21-obs-sweep.md:98](docs/validation/directive-v21-obs-sweep.md#L98)、[directive-v21-obs-sweep.md:242](docs/validation/directive-v21-obs-sweep.md#L242)。

**7. G7：三层审查是否过重**
结论【判断，Medium；P35“有无裁剪文字”为 none】：P35 已有明确逐环裁剪，不能再指控其没有裁剪说明；一般治理片的三层 PASS 则不能当作三份独立实质保证。

- **P35 事实**：落盘定义明确跳过①②③④⑤⑥，⑧a由主控逐字转录，⑦不适用，保留⑨/⑧c/⑩；报告及证据包仍是交付物，所以这里的“无 diff”只能按“无实现 diff”理解，不能扩大为归档也免检。本轮不应再强加⑤⑥⑦。[REMAINING_SLICES.md:1174](docs/t027/REMAINING_SLICES.md#L1174)、[REMAINING_SLICES.md:1179](docs/t027/REMAINING_SLICES.md#L1179)。不过，台账文字不等于已证明同轮授权及全部例外要素齐备。
- **实质失守实例，High**：准则允许特例处理仓库内误放位置，明确排除合法片目录内仅次序违规；P34 ⑥却写成 `position valid`，并声称排除 tmp 根、仓库根和项目子目录，仍判 PASS。这是适用域反向理解，不是翻译小差异；后续终审摘要也未见识别该问题。[DELIVERY_DIRECTIVE.md:516](docs/DELIVERY_DIRECTIVE.md#L516)、[directive-tmp-self-heal.md:140](docs/validation/directive-tmp-self-heal.md#L140)、[directive-tmp-self-heal.md:162](docs/validation/directive-tmp-self-heal.md#L162)、[directive-tmp-self-heal.md:173](docs/validation/directive-tmp-self-heal.md#L173)。
- **不能反向全盘否定**：P31 ⑥抓到成员真值源缺口；其⑦又如实声明收窄 PASS 不等于全片通过。应比较各环新增证据，而非只数合法结论。批量辅助本就不是必经审查层，不能额外算成第四道强制门。[directive-ledger-field-split.md:85](docs/validation/directive-ledger-field-split.md#L85)、[directive-ledger-field-split.md:156](docs/validation/directive-ledger-field-split.md#L156)、[DELIVERY_DIRECTIVE.md:296](docs/DELIVERY_DIRECTIVE.md#L296)。

**8. G8：缺乏执行或判定证据的具体条款**
结论【判断，Medium】：至少以下三项不能凭“报告齐全、门禁全绿”认定已经落实；其中记录缺失与完全未执行必须区分。

- **§7.7 每四片抽查的逐片台账：Medium。** 正文要求自窗口首片起，无抽查也登记一行；唯一登记块仍只有初始化说明及“窗口未满四片”，没有后续逐片记录。可确认的是登记义务缺失，不能据此断言从未实际盲审。[DELIVERY_DIRECTIVE.md:783](docs/DELIVERY_DIRECTIVE.md#L783)、[DELIVERY_DIRECTIVE.md:1157](docs/DELIVERY_DIRECTIVE.md#L1157)。
- **R7.1 报告逐调用六字段：Medium。** 正文要求报告按 `{role, model, stage, startedAt, reason, retryOf}` 逐项列出，P34 成本节却只是角色聚合次数，还把已重发的⑥列为一次。至少报告层没有兑现该要求；原始会话日志是否完整，本轮不知道。[DELIVERY_DIRECTIVE.md:771](docs/DELIVERY_DIRECTIVE.md#L771)、[directive-tmp-self-heal.md:245](docs/validation/directive-tmp-self-heal.md#L245)、[directive-tmp-self-heal.md:171](docs/validation/directive-tmp-self-heal.md#L171)。
- **§12.1–12.2 的有效性判定：Medium。** 五版期间试点顺延从第21次到第25次，所读材料未提供新的通过判定。顺延有裁决记录，不等于越权；但这些治理片的 697/697 也不能替代试点指标，更不能支持“准则已验证”。正文自己保留了这一禁令。[DELIVERY_DIRECTIVE.md:1143](docs/DELIVERY_DIRECTIVE.md#L1143)、[DELIVERY_DIRECTIVE.md:1148](docs/DELIVERY_DIRECTIVE.md#L1148)、[DELIVERY_DIRECTIVE.md:998](docs/DELIVERY_DIRECTIVE.md#L998)。
- **异常登记豁免的实用判定：Low，证据不足。** v1.55 已允许相应触发不阻断提交，但本片样本仍选择码位和同义措辞避开触发，未展示一次按该豁免完成的判定链。这只能说明本样本未证实其收益，不能说全局从未生效；明示留待实现的标记机制不计作现行违约。[DELIVERY_DIRECTIVE.md:673](docs/DELIVERY_DIRECTIVE.md#L673)、[directive-v21-obs-sweep.md:90](docs/validation/directive-v21-obs-sweep.md#L90)、[directive-v21-obs-sweep.md:262](docs/validation/directive-v21-obs-sweep.md#L262)。

**未展开项清单**
历史逐版 diff 与精确条款净增；P35 原始起跑消息及完整同轮授权；逐调用原始日志、真实 token 和当前耗时；其余判据函数及全部夹具的完整审阅；P35 证据包哈希自洽性复验。以上均未据聚合报告推定成立。

本轮自评：事实/判断约 **7:3**，为论证内容估计，非机械统计。
