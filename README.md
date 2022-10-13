#### Переводчик вслух
#### Возможность чтения на двух языках

- Копируешь текст на иностранном языке и программа читает в слух

#### Горячие клавиши

```
ctrl+c - скопировать и прочитать на английском
ctrl+z - повторить на русском
ctrl+alt+p  - Пауза
```

#### Как запустить?

- Требуется golang
- Установить зависимости

```
sudo -H pip3 install gTTS
sudo apt install -y mpg123 translate-shell \
        xsel xclip \
        xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

#### Как запустить?

```
make run
```
