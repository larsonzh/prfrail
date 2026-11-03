# GLM 批量辅助 #1 —— prfrail CLI flag 定义面机械检索（原文归档）

载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`；只读检索；实得输出字符数 7995（自报 `OUTPUT_CHARS=7995`）。

## 1. 命令路由表
`CLI.Execute` switch（74-107；case/定义行，未标文件=cli.go）：help/--help/-h→writeMainUsage 76；version→executeVersion 79/117；init→executeInit 81/133；validate→executeValidate 83/190；preview→executePreview 85/preview.go:37；approvals→executeApprovals 87/approvals.go:57；interactions→executeInteractions 89/interactions.go:46；cost→executeCost 91/cost.go:36；secret→executeSecret 93/secret.go:17；ai→executeAI 95/ai.go:17；run→executeRun 97/274；supervised→executeSupervised 99/supervised_cmd.go:85；report→executeReport 101/544；config(explain)→executeConfigExplain 106/227；default→unknown command+usage 108-110。
二级 case（同父文件，→定义行）：ai（ai.go）：check 23→33；verify 25→116；default 用法 213。approvals（approvals.go）：list 63→74；revoke 65→117。interactions（interactions.go）：list 52→65；respond 54→154；tui 56→98。cost（cost.go）：report 42→51。secret（secret.go）：set 23→35；status 25→68；delete 27→89。supervised（supervised_cmd.go）：return 91→104；report 93→133。run 内部：executeSupervisedRun SU:34,executeAgentRunnerRun 428。

## 2. flag 定义表
文件（internal/console/）：CLI=cli.go PV=preview.go AI=ai.go AP=approvals.go IN=interactions.go CO=cost.go SE=secret.go SU=supervised_cmd.go。图例：A=output as JSON B=path to chain config file C=run directory D=run identifier E=path to authorization ledger file F=path to operator interaction ledger G=configured AI channel H=maximum provider requests。码：sr=supervised return sp=supervised report ac=ai check av=ai verify al=approvals list ar=approvals revoke il=interactions list it=interactions tui ir=interactions respond cr=cost report sd=secret status/delete
chain|string|""|B|init CLI:136,validate CLI:192,preview PV:39,config explain CLI:229,run CLI:276,ac AI:36,av AI:119,sr SU:106
json|bool|false|A|version CLI:119,init CLI:138,validate CLI:193,preview PV:40,config explain CLI:230,run CLI:279,report CLI:547,sr SU:110,sp SU:136,ac AI:42,av AI:132,al AP:77,ar AP:125,il IN:68,ir IN:162,cr CO:54,secret set SE:39,sd SE:112
run-dir|string|""|C|run CLI:277,report CLI:546,sr SU:107,sp SU:135
run-id|string|""|D|run CLI:278,sr SU:108
init|workspace|string|"."|CLI:135|workspace root
init|force|bool|false|CLI:137|overwrite output config file
run|request-record|string|""|CLI:280|path to a declared AgentRunner request record (enables the AgentRunner assembly)
run|capability-record|string|""|CLI:281|path to a declared AgentRunner capability record
run|enforcement-record|string|""|CLI:282|path to a declared AgentRunner enforcement record
run|availability-record|string|""|CLI:283|path to a declared AI availability record
run|authorization-ledger|string|""|CLI:284|path to an authorization ledger
run|cost-ledger|string|""|CLI:285|path to a cost ledger
run|agent-cli|string|""|CLI:286|path to the pinned agent CLI executable
run|agent-version|string|""|CLI:287|pinned agent CLI version
sr|outcome|string|""|SU:109|return outcome: complete or abort
ac|copilot|string|""|AI:38|path to the pinned Copilot CLI executable
ac|workspace|string|""|AI:39|probe working directory
ac|out|string|""|AI:40|write the immutable availability record
channel|string|""|G|ac AI:37,av AI:129
max-requests|int|0|H|ac AI:41,av AI:131
av|record|string|""|AI:120|path to an AI availability record
av|prior-record|Func|无|AI:122|path to a prior AI availability record; repeatable
av|max-age|duration|0|AI:130|maximum record age
ledger|string|""|E|al AP:76,ar AP:121
ar|authorization-id|string|""|AP:119|authorization id to revoke
ar|reason-code|string|""|AP:120|revocation reason code
ar|run-dir|string|""|AP:122|run directory for managed process identity
ar|stop|string|"requested"|AP:123|stop disposition: not-required or requested
ar|actor-id|string|"local-operator"|AP:124|revoking operator identity
actor-id|string|"local-operator"|responding operator identity|it IN:101,ir IN:159
ledger|string|""|F|il IN:67,it IN:100,ir IN:156
ir|interaction-id|string|""|IN:157|interaction id
ir|selection|string|""|IN:158|one allowed response id
ir|attempt|int|0|IN:160|expected task attempt
ir|context-hash|string|""|IN:161|expected context hash
cr|ledger|string|""|CO:53|path to cost ledger file
secret set|ref|string|""|SE:38|opaque secret reference
sd|ref|string|""|SE:111|opaque secret reference（parseSecretReference SE:109）

## 3. 重复标志（multi-value）现状
存在一处：`ai verify` 的 `prior-record`（ai.go:121-127）用标准库 `(*flag.FlagSet).Func` 注册可重复 flag，闭包 append 到 `priorRecordPaths []string`（空值报错）。其余全为单值 Bool/String/Int/Duration；无自定义 flag.Value、flag.Func、set.Visit；`set.Args()` 仅拒绝位置参数（SE:43/117,CO:58），非多值累加。

## 4. run 与 supervised 的装配输入端
executeRun（274）：307-308 FindExecutableSteps/FindSupervisedSteps；309-318 组 inputs；321 `len(supervised)>0 && inputs.anyProvided()`→互斥错误；325 `len(executable)>0 && !inputs.complete() && len(supervised)==0`→报仅支持 noop；357-364 路由：supervised>0→executeSupervisedRun(358)；executable>0→executeAgentRunnerRun(361)；否则 ExecuteNoopRun(363)。
agentRunnerRunInputs（376-385）八字段（均 string）：requestRecord/capabilityRecord/enforcementRecord/availabilityRecord/authorizationLedger/costLedger/agentCLI/agentVersion。判定：complete() 387（TrimSpace 全非空）；anyProvided() 407（任一非空）。
executeAgentRunnerRun（428）→ AgentRunnerAssemblyOptions（504-523）：505 RunID←runID；506 RunRoot←runDir；507 Definition←definition；508 RequestRecord←requestRecord(444-448)；509 CapabilityRecord←capabilityRecord(452-456)；510 EnforcementRecord←&enforcementRecord(460-464)；511 AvailabilityRecord←availabilityRecord(468-472)；512 AvailabilityPolicy←AgentRunnerDefaultAvailabilityPolicy(availabilityRecord,cli.now())；513 AuthorizationLedger←loadAuthorizationLedger(476-480)；514 CostLedger←LoadCostLedger(DecodeCostLedgerRecord(488-496))；515-521 AgentCLI←AgentRunnerCLIOptions{Executable: agentCLIPath(500),Version: inputs.agentVersion(517),VersionArgs: agentRunnerDefaultVersionFlag(518),Timeout: agentRunnerDefaultRunTimeout(519),MaxLogBytes: agentRunnerDefaultMaxLogBytes(520)}；522 Clock←cli.now；524 ExecuteAgentRunnerRun。
supervised 族无等价记录输入端：executeSupervisedReturn/Report 仅 chain/run-dir/run-id/outcome/run-dir 标志；装配输入由 supervisedAssemblyOptions（supervised_cmd.go:50-66）从 definition+runDir 组装（VisibleSender/GateRunner/Clock）⇒ 无记录输入端。

## 5. 用法字符串（usage）落点
（省略共同前缀 usage: prfrail）
ai.go:34 ai check --channel agent-runner-cli --copilot <path> --max-requests 1 --out <path> [--chain <path>] [--wo…
ai.go:117 ai verify --record <path> --channel <channel> --max-age <duration> --max-requests 1 [--prior-record <path…
ai.go:213 ai <check|verify> [options]
approvals.go:59,69 approvals <list|revoke> [options]
approvals.go:79,82 approvals list [--ledger <path>] [--json]
approvals.go:127,130,133,136 approvals revoke --authorization-id <id> --reason-code <code> [--ledger <path>] [--run-dir <path>] [--sto…
cli.go:105,232,235 config explain [--chain <path>] [--json]
cli.go:121,124 version [--json]
cli.go:140,143 init [--workspace <path>] [--chain <path>] [--force] [--json]
cli.go:195,198 validate [--chain <path>] [--json]
cli.go:370 run [--chain <path>] [--run-id <id>] [--run-dir <path>] [--json] [--request-record <path> --capability-rec…
cli.go:549,552,555 report --run-dir <path> [--json]
cli.go:598 <command> [options]
cost.go:38,46,56,59 cost report [--ledger <path>] [--json]
interactions.go:48,60 interactions <list|respond> [options]
interactions.go:70,73 interactions list [--ledger <path>] [--json]
interactions.go:103,106 interactions tui [--ledger <path>] [--actor-id <id>]
interactions.go:163 interactions respond --interaction-id <id> --selection <id> --attempt <n> --context-hash <hash> [--act…
preview.go:42,45 preview [--chain <path>] [--json]
secret.go:19,30 secret <set|status|delete> --ref <reference> [--json]
secret.go:36 secret set --ref <reference> [--json]
secret.go:69 secret status --ref <reference> [--json]
secret.go:90 secret delete --ref <reference> [--json]
supervised_cmd.go:102 supervised <return|report> [options]
supervised_cmd.go:111 supervised return --run-id <id> --outcome complete|abort [--run-dir <path>] [--chain <path>] [--json]
supervised_cmd.go:137 supervised report --run-dir <path> [--json]

OUTPUT_CHARS=7995
