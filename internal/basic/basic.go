package basic

import "context"

type VoiceInterface interface {
	Start(ctx context.Context) error
	Stop()
	SpeedSub() (string, float64, error)
	SpeedAdd() (string, float64, error)
	SpeekLoop(ctx context.Context) error
	Requset(string, string) (string, error)
	Say() error
	ChSpeakMe(string)
	Exit()
	GetPause() bool
	SetPause()
	TanslateOrNot() bool
	InvertTranslate()
}
