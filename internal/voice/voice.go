package voice

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/exitstop/speaker_alpine/internal/google"
	"github.com/exitstop/speaker_alpine/internal/logger"
	"go.uber.org/zap"
)

type VoiceStore struct {
	IP          string
	Port        string
	SpeakMe     google.ChanTranslateMe
	Client      *http.Client
	ChanSpeakMe chan google.ChanTranslateMe
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

	v.ChanSpeakMe = make(chan google.ChanTranslateMe)
	v.Terminatate = make(chan bool)
	v.ChanPause = make(chan bool)

	v.SpeechSpeed = 3.7

	return v
}

func (v *VoiceStore) Start(ctx context.Context) (err error) {
	out, err := v.Requset("get_engine", `{"Text": ""}`)
	if err != nil {
		return err
	}

	logger.Log.Info("gstore.LoopTransalate",
		zap.String("out", string(out)),
	)
	defer logger.Log.Sync()

	// com.google.android.tts com.acapelagroup.android.tts
	out, err = v.Requset("set_engine", `{"Text": "com.google.android.tts"}`)
	if err != nil {
		return err
	}

	logger.Log.Info("set_engine",
		zap.String("out", string(out)),
	)

	out, err = v.Requset("set_speech_rate", `{"SpeechRate": 3}`)
	if err != nil {
		return err
	}

	logger.Log.Info("set_speech_rate",
		zap.String("out", string(out)),
	)

	str := fmt.Sprintf(`{"Text": "Инициализация успешна"}`)
	out, err = v.Requset("play_on_android", str)
	if err != nil {
		return err
	}

	logger.Log.Info("play_on_android",
		zap.String("out", string(out)),
	)

	return v.SpeekLoop(ctx)
}

func (v *VoiceStore) SpeedSub() (out string, speed float64, err error) {
	v.SpeechSpeed -= 0.1
	strSpeed := fmt.Sprintf(`{"SpeechRate": %.2f}`, v.SpeechSpeed)
	out, err = v.Requset("set_speech_rate", strSpeed)
	out = fmt.Sprintf("%s %.1f", out, v.SpeechSpeed)
	speed = v.SpeechSpeed
	return
}

func (v *VoiceStore) SpeedAdd() (out string, speed float64, err error) {
	v.SpeechSpeed += 0.1
	strSpeed := fmt.Sprintf(`{"SpeechRate": %.2f}`, v.SpeechSpeed)
	out, err = v.Requset("set_speech_rate", strSpeed)
	out = fmt.Sprintf("%s %.1f", out, v.SpeechSpeed)
	speed = v.SpeechSpeed
	return
}

func (v *VoiceStore) SpeekLoop(ctx context.Context) (err error) {
	for {
		select {
		case v.SpeakMe = <-v.ChanSpeakMe:
		case <-ctx.Done():
			err = fmt.Errorf("ctx.Done")
			return
		case v.Pause = <-v.ChanPause:
			v.Pause = <-v.ChanPause
			v.SpeakMe.Translate = "пауза снята"
		}

		err = v.Say(ctx, "ru", "")

		if err != nil {
			logger.Log.Info("SpeakLoop",
				zap.String("err", err.Error()),
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

func (v *VoiceStore) Say(ctx context.Context, lang, text string) (err error) {
	str := fmt.Sprintf(`{"Text": "%s"}`, v.SpeakMe)
	out, err := v.Requset("play_on_android", str)

	logger.Log.Info("play_on_android",
		zap.String("out", string(out)),
	)
	logger.Log.Sync()
	return
}

func (v *VoiceStore) ChSpeakMe(in google.ChanTranslateMe) {
	v.ChanSpeakMe <- in
	return
}
func (v *VoiceStore) Exit() {
	fmt.Println("VoiceStore exit")
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
