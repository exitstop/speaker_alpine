package gtts

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/exitstop/speaker_alpine/internal/google"
	"github.com/exitstop/speaker_alpine/internal/logger"
	"go.uber.org/zap"
)

type VoiceStore struct {
	IP              string
	Port            string
	SpeakMe         google.ChanTranslateMe
	Client          *http.Client
	ChanSpeakMe     chan google.ChanTranslateMe
	Terminatate     chan bool
	ChanPause       chan bool
	Pause           bool
	SpeechSpeed     float64
	NoTranslate     bool // не переводить текст
	c1              *exec.Cmd
	c2              *exec.Cmd
	DoubleTranslate bool // двойной перевод
}

func Create() (v VoiceStore) {
	v.Client = &http.Client{
		Timeout: 2 * time.Second,
	}

	v.ChanSpeakMe = make(chan google.ChanTranslateMe)
	v.Terminatate = make(chan bool)
	v.ChanPause = make(chan bool)

	v.SpeechSpeed = 3

	return v
}

func (v *VoiceStore) Start(ctx context.Context) (err error) {
	logger.Log.Info("gtts",
		zap.String("command", "sudo -H pip3 install gTTS; sudo apt install -y mpg123"),
	)
	logger.Log.Sync()

	v.SpeakMe.Translate = "инициализация успешна"
	v.Say(ctx, "ru", v.SpeakMe.Translate)
	v.SpeekLoop(ctx)

	return
}

func (v *VoiceStore) SpeedSub() (out string, speed float64, err error) {
	v.SpeechSpeed -= 1
	speed = v.SpeechSpeed
	return
}

func (v *VoiceStore) SpeedAdd() (out string, speed float64, err error) {
	v.SpeechSpeed += 1
	speed = v.SpeechSpeed
	return
}

func (v *VoiceStore) SpeekLoop(ctx context.Context) (err error) {
	ctxSpeak, cancel := context.WithCancel(ctx)
	for {
		select {
		case v.SpeakMe = <-v.ChanSpeakMe:
		case <-ctx.Done():
			err = fmt.Errorf("terminatate gtts")
			return
		case v.Pause = <-v.ChanPause:
			v.Pause = <-v.ChanPause
			v.SpeakMe.Translate = "пауза снята"
		}

		cancel()
		//v.c1.Process.Kill()
		ctxSpeak, cancel = context.WithCancel(ctx)

		logger.Log.Info("Say",
			zap.String("Translate", v.SpeakMe.Translate),
			zap.String("Orig", v.SpeakMe.Orig),
		)
		logger.Log.Sync()

		go func() {
			err := v.Say(ctxSpeak, "ru", v.SpeakMe.Translate)

			if err != nil {
				logger.Log.Info("Say",
					zap.String("error", err.Error()),
				)
				logger.Log.Sync()
				return
			}
			if v.SpeakMe.Orig != "" && v.DoubleTranslate {
				err = v.Say2(ctxSpeak, "en", v.SpeakMe.Orig, 1)

				if err != nil {
					logger.Log.Info("Say",
						zap.String("error", err.Error()),
					)
					logger.Log.Sync()
					return
				}
			}
		}()

		//time.Sleep(time.Second * 1)
	}
	return
}

func (v *VoiceStore) Stop() {
}

func (v *VoiceStore) Requset(method, input string) (out string, err error) {
	url := fmt.Sprintf("http://%s/%s", v.IP, method)
	data := []byte(input)
	r := bytes.NewReader(data)

	resp, err := v.Client.Post(url, "application/json", r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	out = string(bodyBytes)
	return
}

func (v *VoiceStore) Say(ctx context.Context, lang, text string) (err error) {
	strCommand := fmt.Sprintf(`gtts-cli -l %s "%s"`, lang, text)
	v.c1 = exec.CommandContext(ctx, "/bin/bash", "-c", strCommand)
	stdout1, err := v.c1.StdoutPipe()
	err = v.c1.Start()
	if err != nil {
		return
	}

	strCommand2 := fmt.Sprintf(`mpg123 -d %d --pitch 0 -`, int(v.SpeechSpeed))
	v.c2 = exec.CommandContext(ctx, "/bin/bash", "-c", strCommand2)
	v.c2.Stdin = stdout1
	err = v.c2.Start()

	if err != nil {
		return
	}
	err = v.c1.Wait()
	if err != nil {
		return
	}
	err = v.c2.Wait()
	if err != nil {
		return
	}

	return
}

func (v *VoiceStore) Say2(ctx context.Context, lang, text string, speed int) (err error) {
	strCommand := fmt.Sprintf(`gtts-cli -l %s "%s"`, lang, text)
	v.c1 = exec.CommandContext(ctx, "/bin/bash", "-c", strCommand)
	stdout1, err := v.c1.StdoutPipe()
	err = v.c1.Start()
	if err != nil {
		return
	}

	strCommand2 := fmt.Sprintf(`mpg123 -d %d --pitch 0 -`, speed)
	v.c2 = exec.CommandContext(ctx, "/bin/bash", "-c", strCommand2)
	v.c2.Stdin = stdout1
	err = v.c2.Start()

	if err != nil {
		return
	}
	err = v.c1.Wait()
	if err != nil {
		return
	}
	err = v.c2.Wait()
	if err != nil {
		return
	}

	return
}

func (v *VoiceStore) ChSpeakMe(in google.ChanTranslateMe) {
	v.ChanSpeakMe <- in
	return
}
func (v *VoiceStore) Exit() {
	v.Terminatate <- true
}
func (v *VoiceStore) GetPause() bool {
	return v.Pause
}
func (v *VoiceStore) SetPause() {
	v.Pause = !v.Pause
}
func (v *VoiceStore) TanslateOrNot() (ret bool) {
	return v.NoTranslate
}
func (v *VoiceStore) InvertTranslate() {
	v.NoTranslate = !v.NoTranslate
}
