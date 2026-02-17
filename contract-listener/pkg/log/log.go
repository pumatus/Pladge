package log

import "log"

func Info(v ...any) {
	log.Println(v...)
}

func Error(v ...any) {
	log.Println(v...)
}

func Fatal(v ...any) {
	log.Fatal(v...)
}

func Panic(v ...any) {
	log.Panic(v...)
}
