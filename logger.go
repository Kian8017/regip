package regip

import (
	"fmt"
	"github.com/gookit/color"
	"io"
	"log"
	"time"
)

type Logger struct {
	oup    *log.Logger
	prefix string
}

func NewLogger(where io.Writer) *Logger {
	var l Logger
	l.oup = log.New(where, "", log.LstdFlags)
	l.prefix = ""
	return &l
}

func (l *Logger) Print(what ...interface{}) {
	l.oup.Print(append([]interface{}{l.prefix}, what...)...)
}

func (l *Logger) Printf(format string, what ...interface{}) {
	l.oup.Printf(l.prefix+format, what...)
}

func (l *Logger) Println(what ...interface{}) {
	//l.oup.Println(what...)
	l.oup.Print(append([]interface{}{l.prefix}, what...)...)
}

func (l *Logger) Error(what ...interface{}) {
	l.Print(color.Error.Render(what...))
}

func (l *Logger) Errorf(format string, what ...interface{}) {
	l.Print(color.Error.Render(fmt.Sprintf(format, what...)))
}

func (l *Logger) Time(st time.Time, event string) {
	d := time.Since(st)
	l.Printf("%s took %s", event, CLR_time.Render(d))
}

func (l *Logger) Tag(text string, clr color.Color) *Logger {
	var nl Logger
	nl.oup = l.oup
	nl.prefix = l.prefix + clr.Render("["+text+"] ")
	return &nl
}

func (l *Logger) RawMessage(m *Message) {
	enc, err := m.Marshal()
	l.Tag("RAW", CLR_api).Print(string(enc), " ", err)
}
