package log

func New(logconfig LogConfig) *Logger {
	w := &LogWriter{
		rec:       make(chan *Record, logconfig.BufferLength),
		isEnd:     make(chan bool),
		LogConfig: logconfig,
	}

	w.parameterization()

	if err := w.rotate(); err != nil {
		return nil
	}

	go w.run()

	return &Logger{
		Level:  logconfig.Level,
		Writer: w,
	}
}
