#### Переводчик в слух

- Копируешь текст на иностранном языке и программа читает в слух

#### Как запустить?

- Требуется golang
- Установить зависимости

```
sudo -H pip3 install gTTS
sudo apt install -y mpg123 \
        xsel xclip \
        xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev \

go get github.com/playwright-community/playwright-go
go install github.com/playwright-community/playwright-go/cmd/playwright
playwright install --with-deps
```

- Запускаем, вариант с google tts

```
make google_speech
```
