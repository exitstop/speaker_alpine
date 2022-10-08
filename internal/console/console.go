package console

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/eiannone/keyboard"

	//"github.com/exitstop/robotgo"
	"github.com/exitstop/speaker_alpine/internal/basic"
	"github.com/exitstop/speaker_alpine/internal/google"
	hook "github.com/robotn/gohook"
	"github.com/sirupsen/logrus"
)

type model struct {
	LogLevel    int
	MaxLogLevel int
}

var mod = model{
	LogLevel:    0,
	MaxLogLevel: len(LogLevelString),
}

var LogLevelString = [...]string{
	"info",
	"trace",
	"debug",
	"warning",
	"error",
	"fatal",
}

func (m *model) LevelIntToString() {
	m.LogLevel++
	id := m.LogLevel % m.MaxLogLevel
	if m.LogLevel >= m.MaxLogLevel {
		m.LogLevel = 0
	}

	switch id {
	case 0: //nolint:goconst
		logrus.SetLevel(logrus.InfoLevel)
	case 1: //nolint:goconst
		logrus.SetLevel(logrus.TraceLevel)
	case 2: //nolint:goconst
		logrus.SetLevel(logrus.DebugLevel)
	case 3:
		logrus.SetLevel(logrus.WarnLevel)
	case 4: //nolint:goconst
		logrus.SetLevel(logrus.ErrorLevel)
	case 5: //nolint:goconst
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	fmt.Println("logLevel", LogLevelString[id])
}

func Keyboard() (err error) {
	if err = keyboard.Open(); err != nil {
		return
	}

	defer func() {
		_ = keyboard.Close()
	}()

FOR0:
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyCtrlC:
			break FOR0
		}

		switch char {
		case 'q':
			break FOR0
		case 'c':
			break FOR0
		case 'l':
			mod.LevelIntToString()
		}

		if key == keyboard.KeyEsc {
			break FOR0
		}
	}
	//os.Exit(0)
	return
}

func Add(cancel context.CancelFunc, gstore *google.GStore, voice basic.VoiceInterface) {
	fmt.Println("--- Please press ctrl + q to stop hook ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-q")
		chSpeak := google.ChanTranslateMe{
			Translate: "завершение программы",
		}
		voice.ChSpeakMe(chSpeak)
		time.Sleep(1 * time.Second)
		cancel()
	})

	hook.Register(hook.KeyDown, []string{"p", "ctlr", "alt"}, func(e hook.Event) {
		fmt.Println("ctrl-alt-p")

		if !voice.GetPause() {
			chSpeak := google.ChanTranslateMe{
				Translate: "пауза",
			}
			voice.ChSpeakMe(chSpeak)
		} else {
			chSpeak := google.ChanTranslateMe{
				Translate: "пауза снята",
			}
			voice.ChSpeakMe(chSpeak)
		}

		//time.Sleep(time.Millisecond * 1)
		voice.SetPause()
	})

	hook.Register(hook.KeyDown, []string{"t", "alt"}, func(e hook.Event) {
		fmt.Println("alt-t")

		voice.InvertTranslate()

		if voice.TanslateOrNot() {
			chSpeak := google.ChanTranslateMe{
				Translate: "без перевода",
			}
			voice.ChSpeakMe(chSpeak)
		} else {
			chSpeak := google.ChanTranslateMe{
				Translate: "переводить текст",
			}
			voice.ChSpeakMe(chSpeak)
		}
	})

	hook.Register(hook.KeyDown, []string{"-", "alt"}, func(e hook.Event) {
		fmt.Println("-", "alt")
		out, speed, err := voice.SpeedSub()
		if err != nil {
			fmt.Println(err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"out": out,
		}).Info("speed-")

		str := fmt.Sprintf("%.1f", speed)
		chSpeak := google.ChanTranslateMe{
			Translate: str,
		}
		voice.ChSpeakMe(chSpeak)
	})

	hook.Register(hook.KeyDown, []string{"+", "alt"}, func(e hook.Event) {
		fmt.Println("+", "alt")
		out, speed, err := voice.SpeedAdd()
		if err != nil {
			fmt.Println(err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"out": out,
		}).Info("speed+")

		str := fmt.Sprintf("%.1f", speed)
		chSpeak := google.ChanTranslateMe{
			Translate: str,
		}
		voice.ChSpeakMe(chSpeak)
	})

	fmt.Println("--- Please press c---")
	hook.Register(hook.KeyDown, []string{"c", "ctrl"}, func(e hook.Event) {
		if voice.GetPause() {
			return
		}

		time.Sleep(time.Millisecond * 50)
		text, err := clipboard.ReadAll()

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("clipboard")

			chSpeak := google.ChanTranslateMe{
				Translate: "не скопировалось",
			}
			voice.ChSpeakMe(chSpeak)
			return
		}

		processedString, err := RegexWork(text)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("regexp")
			return
		}

		if voice.TanslateOrNot() {
			processedString, err := RegexWorkRu(text)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Warn("regexp")
				return
			}
			chSpeak := google.ChanTranslateMe{
				Translate: processedString,
				Orig:      text,
			}
			voice.ChSpeakMe(chSpeak)
		} else {
			origText, err := RegexWorkRu(text)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Warn("regexp")
				return
			}

			ch := google.ChanTranslateMe{
				Translate: processedString,
				Orig:      origText,
			}

			select {
			case gstore.ChanTranslateMe <- ch:
				logrus.WithFields(logrus.Fields{
					"SendoToGoole": processedString,
				}).Warn("google")
			default:
				logrus.WithFields(logrus.Fields{
					"SendoToGoole": processedString,
				}).Error("Skip text")
			}
		}
	})

	hook.Register(hook.KeyDown, []string{"r", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("r", "ctrl", "shift")
	})

	s := hook.Start()
	<-hook.Process(s)
}

func Low() {
	EvChan := hook.Start()
	defer hook.End()

	for ev := range EvChan {
		fmt.Println("hook: ", ev)
	}
}

func Event() {
}

func RegexWork(tt string) (out string, err error) {
	tt = strings.ReplaceAll(tt, "\n", ".")
	//tt = strings.ReplaceAll(tt, "»", "")
	//tt = strings.ReplaceAll(tt, "«", "")

	reg0, err := regexp.Compile(`[^a-zA-Z\p{Han}0-9 .,\r\n]+`)
	//reg0, err := regexp.Compile(`[^a-zA-Z0-9 .,]+`)
	if err != nil {
		return
	}
	reg2, err := regexp.Compile(`([\p{L}])\.([\p{L}])`)
	if err != nil {
		return
	}
	reg3, err := regexp.Compile(`([[:lower:]])([[:upper:]])`)
	if err != nil {
		return
	}
	reg4, err := regexp.Compile(`(\b(\p{L}+)\b)`)
	if err != nil {
		return
	}
	tt = reg0.ReplaceAllString(tt, " ")
	tt = reg4.ReplaceAllString(tt, " $1 ")
	tt = reg3.ReplaceAllString(tt, "$1 $2")
	tt = reg2.ReplaceAllString(tt, "$1. $2")

	singleSpacePattern := regexp.MustCompile(`\s+`)
	tt = singleSpacePattern.ReplaceAllString(tt, " ")
	tt = strings.ReplaceAll(tt, " .", ".")
	tt = strings.ReplaceAll(tt, " ,", ",")

	tt = strings.TrimSpace(tt)
	return tt, err
}

func RegexWorkRu(tt string) (out string, err error) {
	reg0, err := regexp.Compile("[^а-яА-Яa-zA-Z0-9 .,]+")
	if err != nil {
		return
	}
	tt = reg0.ReplaceAllString(tt, " ")

	singleSpacePattern := regexp.MustCompile(`\s+`)
	tt = singleSpacePattern.ReplaceAllString(tt, " ")
	tt = strings.ReplaceAll(tt, " .", ".")
	tt = strings.ReplaceAll(tt, " ,", ",")

	tt = strings.TrimSpace(tt)
	return tt, err
}
