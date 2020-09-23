/* Copyright (c) 2020 Kian Musser.
 * This file is part of regip.
 *
 * regip is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * regip is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with regip.  If not, see <https://www.gnu.org/licenses/>.
 */

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
