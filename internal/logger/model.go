package logger

import "go.uber.org/zap"

type LoggerMsg struct {
	Level        string  `json:"level"`
	Microservice string  `json:"microservice"`
	Ts           float64 `json:"ts"`
	Caller       string  `json:"caller"`
	Msg          string  `json:"msg"`
	IdLogger     string  `json:"idLogger"`
	Fields       string  `json:"fields"`
	Error        string  `json:"error"`
}

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции InitLogger, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
type Logger struct {
	*zap.Logger
	originalError error
}

var Log = &Logger{
	Logger: zap.NewNop()}
