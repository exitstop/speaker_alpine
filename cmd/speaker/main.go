package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/exitstop/speaker_alpine/internal/basic"
	"github.com/exitstop/speaker_alpine/internal/console"
	"github.com/exitstop/speaker_alpine/internal/google"
	"github.com/exitstop/speaker_alpine/internal/gtts"
	"github.com/exitstop/speaker_alpine/internal/logger"
	"github.com/exitstop/speaker_alpine/internal/voice"
	"go.uber.org/zap"
)

func main() {
	var (
		nFlag        string
		nGoogleTts   bool
		nFlagTransel bool
		v            basic.VoiceInterface
	)

	flag.BoolVar(&nGoogleTts, "google_speech", false, "use google tts")
	flag.StringVar(&nFlag, "ip", "192.168.0.177", "ip")
	flag.BoolVar(&nFlagTransel, "t", false, "translate")
	flag.Parse()

	logger.Create()

	gstore := google.Create()

	// Запускаем браузер
	if err := gstore.Start(); err != nil {
		log.Println(err)
		return
	}
	defer gstore.Stop()

	if nGoogleTts {
		// google speech
		v_ := gtts.Create()
		v_.DoubleTranslate = true
		gstore.Terminatate = v_.Terminatate
		gstore.SendTranslateToSpeak = v_.ChanSpeakMe
		v = &v_
	} else {
		// voice
		v_ := voice.Create()
		v_.IP = nFlag + ":8484"
		gstore.Terminatate = v_.Terminatate
		gstore.SendTranslateToSpeak = v_.ChanSpeakMe
		v = &v_
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		cancel()
	}()

	go func() {
		//console.Keyboard()
		console.Add(cancel, &gstore, v)
		console.Low()
		//cancel()
	}()

	// Переводчик сам будет слать в chan ChanSpeakMe, чтобы голос воспроизводился
	//go func() {
	//time.Sleep(3 * time.Second)
	//gstore.ChanTranslateMe <- `You can also specify JSHandle as the property value if you want live objects to be passed into the event:`

	//fmt.Println("__0__")
	//time.Sleep(3 * time.Second)
	//gstore.ChanTranslateMe <- `Here you hand errorChannelWatch the errorList as a value.`

	//time.Sleep(3 * time.Second)
	//gstore.ChanTranslateMe <- `To remedy the situation, either hand a slice pointer to errorChannelWatch or rewrite it as a call to a closure, capturing errorList.`
	//}()

	//go func() {
	//for i := 0; i < 20; i++ {
	//time.Sleep(time.Second * 1)
	//str := fmt.Sprintf("Привет мир %d", i)
	//v.ChanSpeakMe <- str
	//}
	//}()

	go func() {
		logger.Log.Info("gstore.LoopTransalate",
			zap.String("init", "ok"),
		)
		logger.Log.Sync()
		// Обработка строк для перевода, посылаемых через ChanTranslateMe
		if err := gstore.LoopTransalate(ctx); err != nil {
			log.Println(err)
			cancel()
			return
		}
	}()

	logger.Log.Info("speaker",
		zap.String("init", "ok"),
	)
	logger.Log.Sync()

	err := v.Start(ctx)
	defer v.Stop()
	if err != nil {
		fmt.Println(err)
		cancel()
		return
	}
}
