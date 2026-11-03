# GLM 批量辅助 1/2：GAP-51/52/55 代码引用面清查原文

载体：runSubagent(agentName="deep-reasoner", model="GLM-5.3-Flash (glm)")。只读事实枚举，逐字归档。

GLM-1/2: DONE

**1. `internal/adapters/bounddir.go`**（共 23 行，纯包级文档注释，无代码）
A/B/C/D：均无（无哨兵、无函数、无错误返回、无测试）。

**2. `internal/adapters/bounddir_unix.go`**（`//go:build !windows && !darwin`，共 111 行）
A：无导出哨兵。
B：
- L22 `type BoundDir struct { fd int; path string }`
- L30 `func OpenBoundDir(path string) (*BoundDir, error)`
- L50 `func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*evidence.BoundFile, error)`
- L56 `func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*evidence.BoundFile, error)`
- L76 `func (dir *BoundDir) Close() error`
- L90 `func boundDirVerify(dir *BoundDir, expected string) error`
C：
- L33 `return nil, &os.PathError{Op: "open", Path: path, Err: err}`
- L38 `return nil, &os.PathError{Op: "stat", Path: path, Err: err}`
- L42 `return nil, fmt.Errorf("adapters: %s is not a directory", path)`
- L58 `return nil, fmt.Errorf("adapters: nil bound directory")`
- L62 `return nil, err`（cleanBoundRelativePath 失败，路径超界守卫）
- L66 `return nil, &os.PathError{Op: "openat", Path: clean, Err: err}`（openat+O_NOFOLLOW，GAP-51 事实：Unix 两段式已在此逐级打开）
- L82 `return err`（Close）
- L92 `return fmt.Errorf("adapters: nil bound directory")`
- L96 `return fmt.Errorf("adapters: stat the bound directory: %v", err)`
- L100 `return err`（os.Stat(expected)）
- L104 `return fmt.Errorf("adapters: %s carries no syscall stat", expected)`
- L107 `return fmt.Errorf("adapters: bound directory %s is not %s", dir.path, expected)`（Dev/Ino 身份比较，路径身份校验）
D：无。

**3. `internal/adapters/bounddir_windows.go`**（`//go:build windows`，共 268 行）
A：无导出哨兵；L37 `var procGetFinalPathNameByHandleW = ...`。
B：
- L22 `const (`（fileFlagBackupSemantics/fileInformationShareAll/volumeNameDos/fileNameNormalized/boundDirNameBufferStart/Limit）
- L40 `type BoundDir struct { handle syscall.Handle; path string }`
- L48 `func OpenBoundDir(path string) (*BoundDir, error)`
- L72 `func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*evidence.BoundFile, error)`
- L78 `func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*evidence.BoundFile, error)`
- L118 `func (dir *BoundDir) confirmOpenedArtifact(handle syscall.Handle, rel string) error`
- L136 `func isUnderDirectory(base, candidate string) bool`（GAP-52 事实：L141-147 用 `strings.EqualFold` 前缀比较）
- L151 `func (dir *BoundDir) Close() error`
- L168 `func boundDirVerify(dir *BoundDir, expected string) error`
- L192 `func (dir *BoundDir) currentName() (string, error)`
- L199 `func handleFinalPath(handle syscall.Handle) (string, error)`
- L220 `func (dir *BoundDir) sameObjectAs(path string) (bool, error)`（GAP-52 事实：L240-242 VolumeSerialNumber+FileIndex 身份比较）
- L243 `func (dir *BoundDir) byHandleFileInformation() (...)`
- L247 `func fileInformationOf(handle syscall.Handle) (...)`
- L257 `func stripExtendedPathPrefix(name string) string`
C：
- L51/L55 `&os.PathError{Op: "open", ...}`；L60 `&os.PathError{Op: "stat", ...}`；L64 `fmt.Errorf("adapters: %s is not a directory", path)`
- L80 `fmt.Errorf("adapters: nil bound directory")`；L84 `return nil, err`（cleanBoundRelativePath 守卫）
- L88 `fmt.Errorf("adapters: resolve the bound directory %s: %v", dir.path, err)`（currentName 失败）
- L93/L97 `&os.PathError{Op: "open", Path: target, ...}`
- L108 `return nil, err`（confirmOpenedArtifact 失败后先 CloseHandle，同源/归属校验拒绝）
- L121 `fmt.Errorf("adapters: resolve the opened artifact %s: %v", rel, err)`
- L125 `fmt.Errorf("adapters: resolve the bound directory %s: %v", dir.path, err)`
- L128 `fmt.Errorf("adapters: artifact %s resolved to %s, which is outside the bound directory %s", rel, final, base)`（isUnderDirectory 判定，前缀比较）
- L157 `return err`（Close）
- L170/L174/L181/L184：boundDirVerify 的 nil/derive/compare/`bound directory %q is not %q` 四类拒绝（L174 走 currentName，L184 走 sameObjectAs 身份比较）
- L209 `return "", callErr`（GetFinalPathNameByHandleW 失败）；L215 `fmt.Errorf("adapters: the handle name exceeds the %d byte buffer", ...)`
- L223/227/231/236 `return false, err`（sameObjectAs 中间失败）；L250 `return syscall.ByHandleFileInformation{}, err`
D：无。

**4. `internal/adapters/bounddir_clean.go`**（`//go:build !darwin`，共 33 行）
A：无。
B：L20 `func cleanBoundRelativePath(rel string) (string, error)`
C：
- L22 `return "", fmt.Errorf("adapters: empty bound artifact name")`
- L25 `return "", fmt.Errorf("adapters: bound artifact name %q must be relative", rel)`（绝对路径拒）
- L29 `return "", fmt.Errorf("adapters: bound artifact name %q escapes the bound directory", rel)`（`..` 超界拒，路径超界守卫）
D：无。

**5. `internal/evidence/boundread.go`**（共 203 行）
A：
- L31 `ErrBoundArtifactExceedsBound = errors.New("bound artifact exceeds its byte bound")`
- L33 `ErrBoundFileClosed = errors.New("bound file is closed")`
- L35 `ErrBoundFileNotRegular = errors.New("bound artifact is not a regular file")`
B：
- L44 `type BoundFile struct { file *os.File; info os.FileInfo; bound int64; closed bool }`
- L55 `func OpenBoundFile(path string, bound int64) (*BoundFile, error)`
- L63 `func OpenBoundFileTruncating(path string, bound int64) (*BoundFile, error)`
- L67 `func openBoundFile(path string, bound int64, allowTruncation bool) (*BoundFile, error)`
- L88 `func NewBoundFile(file *os.File, bound int64, allowTruncation bool) (*BoundFile, error)`
- L122 `func (bound *BoundFile) ReadAll() ([]byte, bool, error)`
- L143 `func (bound *BoundFile) Hash(domain string) (string, error)`
- L165 `func (bound *BoundFile) Identity() (size int64, mod time.Time, mode os.FileMode)`
- L174 `func (bound *BoundFile) Close() error`
- L182 `func (bound *BoundFile) checkReadable() error`
- L196 `func readLimit(bound int64) int64`
C：
- L75 `return nil, err`（openBoundFileHandle 失败，fs.ErrNotExist 透传）
- L90 `fmt.Errorf("%w: an open handle is required", ErrBoundFileNotRegular)`
- L94 `fmt.Errorf("%w: negative byte bound %d", ErrBoundArtifactExceedsBound, bound)`
- L99 `return nil, err`（file.Stat）
- L103 `fmt.Errorf("%w: %s is a directory", ErrBoundFileNotRegular, info.Name())`
- L107 `fmt.Errorf("%w: %s has mode %s", ErrBoundFileNotRegular, info.Name(), info.Mode())`（非 regular 拒，同源/类型身份）
- L111 `fmt.Errorf("%w: %s is %d bytes at the %d byte bound", ...)`（开档超界拒读，!allowTruncation 分支）
- L124 `return nil, false, err`（ReadAll→checkReadable：**关闭后拒读**，ErrBoundFileClosed）
- L127/L131 `return nil, false, err`（Seek/ReadAll 失败）
- L145/L148 `return "", err`（Hash→checkReadable **关闭后拒读** / Seek）
- L154 `return "", err`（io.Copy 失败）
- L157 `fmt.Errorf("%w: digest refused for %s", ErrBoundArtifactExceedsBound, bound.info.Name())`（**超界拒摘要**，written > bound）
- L184 `fmt.Errorf("%w: nil bound file", ErrBoundFileClosed)`
- L187 `fmt.Errorf("%w: %s", ErrBoundFileClosed, bound.info.Name())`（closed 分支）
D：无。

**6. `internal/evidence/boundopen_unix.go`**（`//go:build !windows`，共 14 行）
A/B/C：无哨兵、无拒错返回。B 仅 L11 `func openBoundFileHandle(path string) (*os.File, error)`（L12 `return os.Open(path)`）。

**7. `internal/evidence/boundopen_windows.go`**（`//go:build windows`，共 31 行）
A：无（仅 L15 `const boundFileShareAll = FILE_SHARE_READ|WRITE|DELETE`）。
B：L20 `func openBoundFileHandle(path string) (*os.File, error)`
C：L23 `return nil, &os.PathError{Op: "open", Path: path, Err: err}`（UTF16PtrFromString）；L27 `return nil, &os.PathError{Op: "open", Path: path, Err: err}`（CreateFile，share-delete 模式）。

**8. `internal/evidence/boundread_test.go`**（共 207 行，const `boundFileTestDomain`）
D：
- L33 `TestBoundFileReadsOriginalAfterPathReplacement` — I1：开档后路径被替换（rename+写攻击文件），ReadAll/Hash 仍读原对象。
- L89 `TestBoundFileRefusesOverBoundOnOpen` — I2：record 类超界开档即拒（无 handle 无数据）；恰好 bound 字节可开且不截断。
- L127 `TestBoundFileRefusesDirectory` — I2：目录开档必须失败。
- L141 `TestBoundFileMissingPreservesNotExist` — 缺失工件保持 `fs.ErrNotExist` 可被 errors.Is。
- L158 `TestBoundFileIdentityMatchesRead` — I4：Identity/ReadAll/Hash 同一 handle 一致；bound-1 上限超一字节即拒。
A/B/C：无（测试文件，无哨兵/导出函数/生产拒错路径）。

GLM-1/2: END
