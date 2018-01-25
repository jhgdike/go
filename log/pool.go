package log

import (
	"log"
)

type LogPool struct {
	pool map[string]*Logger
}

var logPool *LogPool

func init() {
	if logPool != nil {
		return
	}
	logPool = &LogPool{pool: make(map[string]*Logger)}
}

// 增加一个日志输出模块
func AddLogger(name, logger *logInterface) {
	logPool.pool[name] = logger
}

func GetLogger(name) (*logInterface, bool) {
	return logPool[name]
}

func (lp *LogPool) log(level int, msg ...interface{}) {
	if len(lp.Pool) < 1 {
		log.Println(msg...)
	} else {
		for _, logger := range lp.Pool {
			if logger.logger != nil &&
				((logger.maxLevel >= level && logger.minLevel <= level) ||
					level == ALL) {
				logger.logAndRotate(msg...)
			}
		}
	}
}

func (lp *LogPool) Get(name string) *Logger {
	return lp.Pool[name]
}

func (lp *LogPool) Debug(msg ...interface{}) {
	baseLog(DEBUG, "[DEBUG]", msg...)
	Callstack(DEBUG)
}

func (lp *LogPool) Info(msg ...interface{}) {
	baseLog(INFO, "[INFO]", msg...)
}

func (lp *LogPool) Warn(msg ...interface{}) {
	baseLog(WARN, "[WARN]", msg...)
	Callstack(WARN)
}

func (lp *LogPool) Error(msg ...interface{}) {
	baseLog(ERROR, "[ERROR]", msg...)
	Callstack(ERROR)
}

func (lp *LogPool) Fatal(msg ...interface{}) {
	baseLog(FATAL, "[FATAL]", msg...)
	Callstack(FATAL)
	os.Exit(1)
}

func baseLog(level int, prefix string, msg ...interface{}) {
	msg = append(logInfo{prefix}, msg...)
	if Log != nil {
		Log.log(level, msg...)
	}
}

func Debug(msg ...interface{}) {
	baseLog(DEBUG, "[DEBUG]", msg...)
	Callstack(DEBUG)
}

func Info(msg ...interface{}) {
	baseLog(INFO, "[INFO]", msg...)
}

func Warn(msg ...interface{}) {
	baseLog(WARN, "[WARN]", msg...)
	Callstack(WARN)
}

func Error(msg ...interface{}) {
	baseLog(ERROR, "[ERROR]", msg...)
	Callstack(ERROR)
}

func Fatal(msg ...interface{}) {
	baseLog(FATAL, "[FATAL]", msg...)
	Callstack(FATAL)
	os.Exit(1)
}

func Callstack(level int) {
	if LogCallstack {
		msg := getCallstack()
		baseLog(level, "", msg...)
	}
}
