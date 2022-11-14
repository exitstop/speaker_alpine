package translateshell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Store struct {
	ctx           context.Context
	ctxSpeak      context.Context
	cancelSpeak   context.CancelFunc
	chText        chan string
	typeOperation string
	pause         bool
	original      string
	translate     string
	lastText      string
}

func New(ctx context.Context) (store *Store) {
	store = &Store{
		ctx:    ctx,
		chText: make(chan string),
	}
	return
}

func (s *Store) Run() {
	s.ctxSpeak, s.cancelSpeak = context.WithCancel(s.ctx)
	for {
		var text string
		select {
		case <-s.ctx.Done():
			return
		case text = <-s.chText:
		}

		s.cancelSpeak()
		time.Sleep(20 * time.Millisecond)
		s.ctxSpeak, s.cancelSpeak = context.WithCancel(s.ctx)

		go func() {
			if text != s.lastText || s.lastText == "" {
				text = strings.ToLower(text)
				s.translate = speak(s.ctxSpeak, text, `trans -b -t ru "%s"`)
				s.lastText = text
				s.original = text
			}

			switch s.typeOperation {
			case operationOnlyTranslate:
				replay(s.ctxSpeak, "ru", s.translate, 5, 2)
			case operationOnlyOriginalRu:
				replay(s.ctxSpeak, "ru", s.original, 5, 2)
			case operationOnlyOriginal:
				replay(s.ctxSpeak, "en", s.original, 2, 1)
			case operationTranslateAndOriginal:
				replay(s.ctxSpeak, "ru", s.translate, 5, 2)
				replay(s.ctxSpeak, "en", s.original, 5, 2)
				//speak(s.ctxSpeak, text, `trans -b -t ru -no-translate -sp "%s"`)
				//default:
				//s.translate = speak(s.ctxSpeak, text, `trans -b -t ru -p "%s"`)
				//speak(s.ctxSpeak, text, `trans -b -t ru -no-translate -sp "%s"`)
			}
		}()
	}
}

const (
	operationOnlyTranslate        string = "OnlyTranslate"
	operationOnlyOriginal         string = "OnlyOriginal"
	operationOnlyOriginalRu       string = "OnlyOriginalRu"
	operationTranslateAndOriginal string = "TranslateAndOriginal"
)

func (s *Store) OnlyTranslate() {
	s.typeOperation = operationOnlyTranslate
}

func (s *Store) OnlyOriginal() {
	s.typeOperation = operationOnlyOriginal
}

func (s *Store) OnlyOriginalRu() {
	s.typeOperation = operationOnlyOriginalRu
}

func (s *Store) TranslateAndOriginal() {
	s.typeOperation = operationTranslateAndOriginal
}

func (s *Store) Go(text string) {
	s.chText <- text
}

func (s *Store) CheckPause() bool {
	return s.pause
}
func (s *Store) SetPause() {
	s.pause = !s.pause
	if s.pause {
		s.cancelSpeak()
	}
}

func speak(ctx context.Context, text, command string) string {
	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf(command, text))
	cmd.Stderr = os.Stderr
	out, _ := cmd.Output()
	return string(out)
}

func replay(ctx context.Context, lang, text string, speed, half int) (err error) {
	strCommand := fmt.Sprintf(`gtts-cli -l %s "%s"`, lang, text)
	fmt.Println(strCommand)
	c1 := exec.CommandContext(ctx, "/bin/bash", "-c", strCommand)
	c1.Stderr = os.Stderr
	stdout1, err := c1.StdoutPipe()
	err = c1.Start()
	if err != nil {
		return
	}

	strCommand2 := fmt.Sprintf(`mpg123 -d %d -h %d --pitch 0 -`, speed, half)
	c2 := exec.CommandContext(ctx, "/bin/bash", "-c", strCommand2)
	c2.Stdin = stdout1
	c2.Stderr = os.Stderr
	err = c2.Start()
	if err != nil {
		return
	}
	err = c1.Wait()
	if err != nil {
		return
	}
	err = c2.Wait()
	if err != nil {
		return
	}

	return
}
