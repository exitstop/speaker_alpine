package google

import (
	"time"

	"github.com/playwright-community/playwright-go"
)

type GStore struct {
	Url                    string
	ToTranslete            string
	LastTranslete          string
	TranslatedText         string
	ChanTranslateMe        chan string
	Drop                   chan struct{}
	Terminatate            chan bool
	SendTranslateToSpeak   chan string
	TimeoutWaitTranslate   time.Duration
	CountLoopWaitTranslate int
	Pw                     *playwright.Playwright
	Browser                playwright.Browser
	Page                   playwright.Page
}
