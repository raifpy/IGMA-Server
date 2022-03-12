package logger

import (
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	//Log *log.Logger
	Log        *zap.Logger
	LogDirPath string
}

func NewLog(LogDirPath string) (l *Log, err error) { //!! Olması gerektiği gibi çalışmıyor. Üzerinde fazla durmadım
	l = &Log{
		LogDirPath: LogDirPath,
	}
	f, err := os.OpenFile(path.Join(LogDirPath, "raw.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return l, err
	}

	l.Log, err = zap.NewProduction(zap.ErrorOutput(zapcore.AddSync(f))) //!!Çalışmıyor

	return

}
