# Chgo

Access Coursehunters content from the comfort of your terminal. You need
to have premium account for this.

[![asciicast](https://asciinema.org/a/nhglYGffSjGg1fDSXYjIHGVRG.svg)](https://asciinema.org/a/nhglYGffSjGg1fDSXYjIHGVRG)

## Installation

```bash
go install github.com/zshbunni/chgo@latest
```

To build from source:

```bash
git clone https://github.com/zshbunni/chgo
cd chgo
go build .
```

## Usage

You need to first login using your Coursehunters email and password to access the TUI.

> NOTE: To play the lessons, chgo uses `mpv`. So make sure that you have installed it.

Login:

```bash
chgo login -u=email -p=pass
```

Once logged in, type `chgo`

## Navigation

| Key Bindings | Description                                                          |
| ------------ | -------------------------------------------------------------------- |
| ctrl+c       | quit                                                                 |
| tab          | toggle focus between search input and courses list                   |
| enter        | fetch lessons for the selected course or play the selected lesson    |
| shift+p      | go back to home screen (this is used when you are in lessons screen) |
