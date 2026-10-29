读取范围：读取主控落盘的 [docs/DEV_PLAN.md](docs/DEV_PLAN.md)、[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md)、[docs/ADR_REGISTER.md](docs/ADR_REGISTER.md)、[docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md)、[docs/validation/directive-d-mechanical-guard.md](docs/validation/directive-d-mechanical-guard.md)、[docs/validation/directive-v21-obs-sweep.md](docs/validation/directive-v21-obs-sweep.md)、[docs/validation/directive-tmp-self-heal.md](docs/validation/directive-tmp-self-heal.md)、[docs/validation/at24-wording.md](docs/validation/at24-wording.md)、[docs/validation/directive-ledger-discipline.md](docs/validation/directive-ledger-discipline.md)、[docs/validation/directive-scan-template.md](docs/validation/directive-scan-template.md)、[docs/validation/directive-ob-sweep.md](docs/validation/directive-ob-sweep.md)、[docs/validation/directive-temp-auth-limit.md](docs/validation/directive-temp-auth-limit.md)。最后五份主要核对头部；其余按相关章节读取。

本轮只读，无文件变更，未执行命令、测试或 lint。以下“事实”指落盘记录事实；统计为人工去重、核对与重算，不代表重新验证实现。

**1. M1 节奏与派生**
结论【判断，Medium】：P 编号窗口内，交付数量主要反映治理吞吐，不能直接代表产品推进；治理工作确实产生了后续治理工作。[REMAINING_SLICES.md:494](docs/t027/REMAINING_SLICES.md#L494)、[DEV_PLAN.md:151](docs/DEV_PLAN.md#L151)

事实：P1–P35 共35个编号块；扣除“不新建片”的P3和待评估的P5–P7，得到31项正式工作，其中30项完成、P35进行中。P3与P19属于同一个C1，不应重复计片；日期窗口内另有未挂P号的 `GATES-EDGE-COVERAGE`，未混入上述分母。[REMAINING_SLICES.md:548](docs/t027/REMAINING_SLICES.md#L548)、[REMAINING_SLICES.md:598](docs/t027/REMAINING_SLICES.md#L598)、[REMAINING_SLICES.md:1356](docs/t027/REMAINING_SLICES.md#L1356)

事实：31项规模为S 7、M 18、M-L 3、L 3。P1日期为10月14日，P34于10月28日完成，P35实际10月29日起跑，不能按名称中的“10-30”计时。30项完成覆盖15个日历日，约每日2项，仅是日历吞吐量。[at24-wording.md:3](docs/validation/at24-wording.md#L3)、[REMAINING_SLICES.md:1374](docs/t027/REMAINING_SLICES.md#L1374)、[REMAINING_SLICES.md:1162](docs/t027/REMAINING_SLICES.md#L1162)

事实：按目标分类，已完成30项中，治理24项，占80%；产品、契约与平台验证5项，另有明确标为非治理的门禁缓存P27。P21–P34这14项集中在10月22日至28日，其中13项治理、1项缓存。[REMAINING_SLICES.md:798](docs/t027/REMAINING_SLICES.md#L798)、[REMAINING_SLICES.md:882](docs/t027/REMAINING_SLICES.md#L882)、[REMAINING_SLICES.md:999](docs/t027/REMAINING_SLICES.md#L999)

事实与判断：至少四条可追溯派生链包含7项已完成后续工作：P2→P15、P13→P14、P27→P28→P29、P30→P32→P33→P34；此外尚有 `GATES-EXEMPT` 和v23候选。不能据此认定每次拆片无价值，但治理存在持续生成新增工作的反馈链，尚未显示稳定收敛。[REMAINING_SLICES.md:544](docs/t027/REMAINING_SLICES.md#L544)、[REMAINING_SLICES.md:736](docs/t027/REMAINING_SLICES.md#L736)、[REMAINING_SLICES.md:1011](docs/t027/REMAINING_SLICES.md#L1011)、[REMAINING_SLICES.md:1096](docs/t027/REMAINING_SLICES.md#L1096)、[directive-tmp-self-heal.md:267](docs/validation/directive-tmp-self-heal.md#L267)

**2. M2 观察项积压**
结论【判断，Medium】：当前状态列夸大了尚需修复的工作量，至少四条应销账或合并引用，而不是再立相同修复片。[DELIVERY_DIRECTIVE.md:1072](docs/DELIVERY_DIRECTIVE.md#L1072)、[REMAINING_SLICES.md:818](docs/t027/REMAINING_SLICES.md#L818)

事实：OB-1至OB-128共128条；按当前处置列统计，“观察中”12条、“待议”11条，合计23条，约18%。观察中为1、7、10、33、87、99、100、104、105、107、111、113；待议为19、21、52、53、55、57、58、59、61、75、120。不含标为“待办”的OB-69。最早OB-1、OB-7登记于9月21日，截至10月29日已38天。[DELIVERY_DIRECTIVE.md:1009](docs/DELIVERY_DIRECTIVE.md#L1009)、[DELIVERY_DIRECTIVE.md:1013](docs/DELIVERY_DIRECTIVE.md#L1013)、[DELIVERY_DIRECTIVE.md:1148](docs/DELIVERY_DIRECTIVE.md#L1148)

事实：OB-52、OB-59仍“待议”，但P17已记完成；OB-57仍“待议”，P18已记根因闭合；OB-58仍“待议”，其后继OB-109已经P20处置。这四条是明确的状态回写滞后。[DELIVERY_DIRECTIVE.md:1077](docs/DELIVERY_DIRECTIVE.md#L1077)、[DELIVERY_DIRECTIVE.md:1129](docs/DELIVERY_DIRECTIVE.md#L1129)、[REMAINING_SLICES.md:818](docs/t027/REMAINING_SLICES.md#L818)、[REMAINING_SLICES.md:841](docs/t027/REMAINING_SLICES.md#L841)

判断：OB-7是正向历史观察，本来不是待修缺陷；OB-111已选择“维持现状”，却仍与待评估项混列。反之，OB-19、OB-55、OB-61、OB-75等仍有明确未来触发条件，不能统一判为无效积压。[DELIVERY_DIRECTIVE.md:1019](docs/DELIVERY_DIRECTIVE.md#L1019)、[REMAINING_SLICES.md:1364](docs/t027/REMAINING_SLICES.md#L1364)、[DELIVERY_DIRECTIVE.md:1075](docs/DELIVERY_DIRECTIVE.md#L1075)、[DELIVERY_DIRECTIVE.md:1095](docs/DELIVERY_DIRECTIVE.md#L1095)

**3. M3 分级使用**
结论【判断，none，针对分级规避指控】：已读材料不足以证明通过分级逃避审查；正式台账几乎全为SLICE，首先是登记规则和治理任务构成造成的选择偏差。[DELIVERY_DIRECTIVE.md:260](docs/DELIVERY_DIRECTIVE.md#L260)

事实：上述31项正式工作中，SLICE 30项、PILOT 1项、ROUTINE 0项，分别约96.8%、3.2%、0%；只看30项已完成工作则为29、1、0。纯ROUTINE按准则不进入该台账，而DEV_PLAN确实记有ROUTINE追加与引用校正，因此不能把台账中的0解释为全仓从未使用ROUTINE。[REMAINING_SLICES.md:839](docs/t027/REMAINING_SLICES.md#L839)、[DELIVERY_DIRECTIVE.md:260](docs/DELIVERY_DIRECTIVE.md#L260)、[DEV_PLAN.md:145](docs/DEV_PLAN.md#L145)

事实与判断【Medium，试点成效】：P18虽然完成，但自捕获率仅1/10，未达到目标75%，材料已明确交用户裁决；“完成PILOT”不等于“试点验证成功”。P35的流程跳过也已明确登记为裁剪，不能作为暗中跳审证据。[REMAINING_SLICES.md:841](docs/t027/REMAINING_SLICES.md#L841)、[REMAINING_SLICES.md:1179](docs/t027/REMAINING_SLICES.md#L1179)

**4. M4 台账维护负担**
结论【判断，Medium】：多位置维护已经产生可观察的不一致，增加登记动作并未相应保证完成记录完整。[REMAINING_SLICES.md:1316](docs/t027/REMAINING_SLICES.md#L1316)、[DEV_PLAN.md:153](docs/DEV_PLAN.md#L153)

事实：台账按末条非空行计1377行；含55个定义或登记块，即历史A/B/C 15块、P编号35块、回写定义5块。它们不是55个独立正式切片。顶层状态表31条，完成记录48条数据行，不含两个表头；准则另有128条OB。[REMAINING_SLICES.md:33](docs/t027/REMAINING_SLICES.md#L33)、[REMAINING_SLICES.md:107](docs/t027/REMAINING_SLICES.md#L107)、[REMAINING_SLICES.md:1218](docs/t027/REMAINING_SLICES.md#L1218)、[REMAINING_SLICES.md:1377](docs/t027/REMAINING_SLICES.md#L1377)

事实：P编号口径的30项已完成工作，只有18项出现在完成记录表；缺P1、P2、P4、P8、P9、P10、P15、P16、P17、P18、P28、P29共12项。这里缺的是完成行，不是据此否定其已有完成证据。[REMAINING_SLICES.md:1316](docs/t027/REMAINING_SLICES.md#L1316)

事实与判断：新片常需维护顶表、定义块、完成行三处并同步双语，另写报告；治理片还维护版本记录、附录C和受影响OB。并非每片都必须同时改DEV_PLAN与ADR。P34的新完成行一面称“只记报告路径”，一面复制三个提交及三次CI，属于既有简化规则未被落实，不需要再新增制度才能纠正。[REMAINING_SLICES.md:1108](docs/t027/REMAINING_SLICES.md#L1108)、[DELIVERY_DIRECTIVE.md:947](docs/DELIVERY_DIRECTIVE.md#L947)、[REMAINING_SLICES.md:1374](docs/t027/REMAINING_SLICES.md#L1374)

**5. M5 分工与重复状态**
结论【判断，Medium】：权威分工基本写清，主要缺口是跨文档更新不能可靠保持一致，而非缺少另一份职责说明。[DELIVERY_DIRECTIVE.md:936](docs/DELIVERY_DIRECTIVE.md#L936)、[DEV_PLAN.md:153](docs/DEV_PLAN.md#L153)

事实：DEV_PLAN负责产品任务、需求映射及阻断；执行台账负责切片定义与进度；验证报告承载证据；治理片以准则附录C登记，不要求回写DEV_PLAN。因此，不能把治理片未进入DEV_PLAN一概判为漏登记。[DELIVERY_DIRECTIVE.md:941](docs/DELIVERY_DIRECTIVE.md#L941)、[DELIVERY_DIRECTIVE.md:951](docs/DELIVERY_DIRECTIVE.md#L951)

事实：已有三种具体漂移：P18/P19/P20完成后漏回写DEV_PLAN，补录说明明确承认门禁看不到该偏离；P3的C1分级与探针仍“待用户确认”，P19同字段已“用户已确认”；P33状态仍写提交、推送、CI待授权，报告后记却已记录完成。[DEV_PLAN.md:153](docs/DEV_PLAN.md#L153)、[REMAINING_SLICES.md:565](docs/t027/REMAINING_SLICES.md#L565)、[REMAINING_SLICES.md:860](docs/t027/REMAINING_SLICES.md#L860)、[REMAINING_SLICES.md:1122](docs/t027/REMAINING_SLICES.md#L1122)、[directive-v21-obs-sweep.md:246](docs/validation/directive-v21-obs-sweep.md#L246)

**6. M6 R/D重算**
结论【判断，High】：P32的R至少为7，不是最终报告所写的5；累计漏计使强制暂停的触发时点失真。[DELIVERY_DIRECTIVE.md:880](docs/DELIVERY_DIRECTIVE.md#L880)、[directive-d-mechanical-guard.md:75](docs/validation/directive-d-mechanical-guard.md#L75)

事实与重算：异常章先明确记录前期R=2；随后⑨三项不同根因各加1；⑧b失败修复再加1；最后一次⑨复跑即使接受报告将两项发现合计1次的宽口径，累计也为7。报告在⑨阶段把局部子计数3当成了全片累计，遗漏前期2次。[directive-d-mechanical-guard.md:61](docs/validation/directive-d-mechanical-guard.md#L61)、[directive-d-mechanical-guard.md:67](docs/validation/directive-d-mechanical-guard.md#L67)、[directive-d-mechanical-guard.md:68](docs/validation/directive-d-mechanical-guard.md#L68)、[directive-d-mechanical-guard.md:69](docs/validation/directive-d-mechanical-guard.md#L69)、[directive-d-mechanical-guard.md:71](docs/validation/directive-d-mechanical-guard.md#L71)、[directive-d-mechanical-guard.md:75](docs/validation/directive-d-mechanical-guard.md#L75)

判断：正确累计在⑨第二项后已经达到4，应触发暂停人工审核，而不是等⑧b。末次两发现若按不同根因拆开，R还可能更高；这不影响“至少7”的结论。D独立累计为3，不应再机械加进R；报告清洁章仍写D=0，也未与最终异常章对齐。[DELIVERY_DIRECTIVE.md:880](docs/DELIVERY_DIRECTIVE.md#L880)、[directive-d-mechanical-guard.md:68](docs/validation/directive-d-mechanical-guard.md#L68)、[directive-d-mechanical-guard.md:77](docs/validation/directive-d-mechanical-guard.md#L77)、[directive-d-mechanical-guard.md:171](docs/validation/directive-d-mechanical-guard.md#L171)

管理含义【判断，High】：R/D仍揭示了登记、自查、暂存及取证方面的反复失败，但当前累计不能可靠承担止损阈值。材料明确登记了用户允许继续的例外，不能称为秘密绕过；也不能用后续授权补偿前面失真的计量。应先校正唯一真值源，再据此讨论阈值是否合理。[directive-d-mechanical-guard.md:71](docs/validation/directive-d-mechanical-guard.md#L71)、[directive-d-mechanical-guard.md:75](docs/validation/directive-d-mechanical-guard.md#L75)、[directive-d-mechanical-guard.md:77](docs/validation/directive-d-mechanical-guard.md#L77)

**7. M7 会话恢复**
结论【判断，Medium】：仓库已有结构化最小阅读入口，但它是规则导航，不是完整的当前会话状态快照。[DELIVERY_DIRECTIVE.md:1349](docs/DELIVERY_DIRECTIVE.md#L1349)、[DELIVERY_DIRECTIVE.md:1406](docs/DELIVERY_DIRECTIVE.md#L1406)

事实：B.7启动卡共69行，含起跑、开场、停点模板和最小阅读集；“所有”行展开后仍包含30个正文小节、B.0及启用角色模板，中等或高风险片再扩读。准则明确禁止每片全文重读，不能把文件总长度说成每次强制阅读量。[DELIVERY_DIRECTIVE.md:1406](docs/DELIVERY_DIRECTIVE.md#L1406)

事实：四份中文核心文档按末条非空行计，准则2158行、台账1377行、DEV_PLAN 213行、ADR 99行，合计3847行；再加最新P34报告291行，共4138行。这是可跳转材料体量，不是实际恢复必读量。[DELIVERY_DIRECTIVE.md:2158](docs/DELIVERY_DIRECTIVE.md#L2158)、[REMAINING_SLICES.md:1377](docs/t027/REMAINING_SLICES.md#L1377)、[DEV_PLAN.md:213](docs/DEV_PLAN.md#L213)、[ADR_REGISTER.md:99](docs/ADR_REGISTER.md#L99)、[directive-tmp-self-heal.md:291](docs/validation/directive-tmp-self-heal.md#L291)

判断：实际恢复至少需取得当前片定义、最近有效停点、产品阻断、适用ADR/OB及按类型选取的规则章节，并确认当轮授权。当前P34报告计数表仍R=2/D=1，台账终态已R=3/D=2，说明仅找到报告并不足以恢复正确游标。先统一现有停点与异常章，比再增加一份重复状态表更直接。[directive-tmp-self-heal.md:96](docs/validation/directive-tmp-self-heal.md#L96)、[REMAINING_SLICES.md:1145](docs/t027/REMAINING_SLICES.md#L1145)、[DELIVERY_DIRECTIVE.md:913](docs/DELIVERY_DIRECTIVE.md#L913)、[DELIVERY_DIRECTIVE.md:1392](docs/DELIVERY_DIRECTIVE.md#L1392)

**未展开项清单**
未编号执行片的全量重建；全仓ROUTINE事件分母及实际工时、费用占比；P32末次两发现的根因拆分；不同角色与切片类型恢复集的逐项展开行数。

本轮自评：事实/判断比例约为75%/25%。
