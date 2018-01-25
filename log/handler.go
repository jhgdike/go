package log

// 根据日志大小切割轮转
func NewSizeRotateLogger(name, path string, maxNumber, maxLevel, minLevel int, maxSize int64) *logInterface {
	Info("add logger", name, path)
	logger := &Logger{
		path:     path,
		filename: name,
		logFile:  openFile(path, name),

		mode:      SIZE_MODE,
		maxNumber: maxNumber,
		maxSize:   maxSize,

		maxLevel: maxLevel,
		minLevel: minLevel,

		suffix: 0,
		mux:    new(sync.RWMutex),
	}
	logger.logger = log.New(logger.logFile, "", log.Ldate|log.Ltime)

	go logger.ListenToReopenFile()
}

// 根据日期切割轮转
func NewDateRotateLogger(name, path string, maxLevel, minLevel int) *logInterface {
	logger := &Logger{
		path:     path,
		filename: name,
		logFile:  openFile(path, name),

		mode:          DATE_MODE,
		lastWriteTime: time.Now(),

		maxLevel: maxLevel,
		minLevel: minLevel,

		suffix: 0,
		mux:    new(sync.RWMutex),
	}
	logger.logger = log.New(logger.logFile, "", log.Ldate|log.Ltime)

	go logger.ListenToReopenFile()
	return logger
}
