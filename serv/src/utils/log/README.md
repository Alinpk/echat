# log
>原生log库，无需额外安装，支持自定义前缀，支持配置输出文件，支持三种记录方式(`println`, `panic`, `fatal`)

# zap
>why: uber开发的日志库，性能好，结构化日志，分级日志
```go
import "go.uber.org/zap"

func main() {
    log := zap.NewExample()
    log.Debug("this is debug message")
    log.Info("this is info message")
    log.Info("this is info message with fields")
}
```

# slog
>go团队开发的结构化日志
>1.支持k-v结构，方便分析搜索
```go
import "log/slog"
textLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
```
>暂时没看关于日志满了的后果
>但是我们可以重新实现fs去去进行日志切割轮转，当然也可以依赖于三方库lumberjack
```go
log := &lumberjack.Logger{
    Filename:   "/path/file.log", // 日志文件的位置
    MaxSize:    10,               // 文件最大尺寸（以MB为单位）
    MaxBackups: 3,                // 保留的最大旧文件数量
    MaxAge:     28,               // 保留旧文件的最大天数
    Compress:   true,             // 是否压缩/归档旧文件
    LocalTime:  true,             // 使用本地时间创建时间戳
}
textLogger := slog.New(slog.NewTextHandler(log, nil))
jsonLogger := slog.New(slog.NewJSONHandler(log, nil))
```
