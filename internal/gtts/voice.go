package gtts

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/exitstop/speaker_alpine/internal/logger"
	"go.uber.org/zap"
)

type VoiceStore struct {
	IP          string
	Port        string
	SpeakMe     string
	Client      *http.Client
	ChanSpeakMe chan string
	Terminatate chan bool
	ChanPause   chan bool
	Pause       bool
	SpeechSpeed float64
	NoTranslate bool // не переводить текст
}

func Create() (v VoiceStore) {
	v.Client = &http.Client{
		Timeout: 2 * time.Second,
	}

	v.ChanSpeakMe = make(chan string)
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

	v.SpeakMe = "инициализация успешна"
	v.Say()
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
	for {
		select {
		case v.SpeakMe = <-v.ChanSpeakMe:
		case <-ctx.Done():
			err = fmt.Errorf("terminatate gtts")
			return
		case v.Pause = <-v.ChanPause:
			v.Pause = <-v.ChanPause
			v.SpeakMe = "пауза снята"
		}

		logger.Log.Info("Say",
			zap.String("text", v.SpeakMe),
		)
		logger.Log.Sync()

		err := v.Say()

		if err != nil {
			logger.Log.Info("Say",
				zap.String("error", err.Error()),
			)
			logger.Log.Sync()
			continue
		}

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

func (v *VoiceStore) Say() (err error) {
	strCommand := fmt.Sprintf(`gtts-cli -l ru "%s" | mpg123 -d %d --pitch 0 -`, v.SpeakMe, int(v.SpeechSpeed))
	cmdCurl := exec.Command("/bin/bash", "-c", strCommand)
	err = cmdCurl.Run()
	return
}

func (v *VoiceStore) ChSpeakMe(in string) {
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
