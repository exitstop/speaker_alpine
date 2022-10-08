package basic

import (
	"context"

	"github.com/exitstop/speaker_alpine/internal/google"
)

type VoiceInterface interface {
	Start(ctx context.Context) error
	Stop()
	SpeedSub() (string, float64, error)
	SpeedAdd() (string, float64, error)
	SpeekLoop(ctx context.Context) error
	Requset(string, string) (string, error)
	Say(ctx context.Context, lang, text string) error
	ChSpeakMe(google.ChanTranslateMe)
	Exit()
	GetPause() bool
	SetPause()
	TanslateOrNot() bool
	InvertTranslate()
}
