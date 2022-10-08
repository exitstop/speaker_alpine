package google

import (
	"time"

	"github.com/playwright-community/playwright-go"
)

type GStore struct {
	Url                    string
	ToTranslete            ChanTranslateMe
	LastTranslete          ChanTranslateMe
	TranslatedText         ChanTranslateMe
	ChanTranslateMe        chan ChanTranslateMe
	Drop                   chan struct{}
	Terminatate            chan bool
	SendTranslateToSpeak   chan ChanTranslateMe
	TimeoutWaitTranslate   time.Duration
	CountLoopWaitTranslate int
	Pw                     *playwright.Playwright
	Browser                playwright.Browser
	Page                   playwright.Page
}

type ChanTranslateMe struct {
	Translate string
	Orig      string
}
