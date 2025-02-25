

# OpenSAI <img src="logo.png" width=40 style="position: absolute; bottom: 0">


Open-source implementation of SAI (SimpleAI) in Go.

![Static Badge](https://img.shields.io/badge/language-Go-blue) ![](https://img.shields.io/github/commit-activity/w/z3nnix/openSAI/main
) ![](https://img.shields.io/github/stars/z3nnix/openSAI?style=flat-square
)

# Setup

First, you need to install dependencies.

### In Debian-based (Ubuntu, etc.)
```bash
sudo apt install golang
```

### In Arch-based (Manjaro, Artix, etc.)
```bash
sudo pacman -S go
```

### In Gentoo-based
```bash
emerge --ask dev-lang/go
```

### Other
Visit https://go.dev/doc/install

## Post install

Then, you need to make neccessary configuration files.
Run commands below in the cloned repo folder.

```bash
cd config && touch info.bot names.bot response.bot token.bot vocabulary.bot && cd ..
```

Now you can build your bot within command below
```bash
go build -o bot cmd/bot/main.go cmd/bot/formatting.go cmd/bot/vocman.go cmd/bot/config.go cmd/bot/fetchHandler.go cmd/bot/amnesiaHandler.go cmd/bot/<ENGINE NAME>.go
```

### Engine list
- StupidV1: A classic SimpleAI engine
- EmbeddingsV1: modified classic engine 

Run it within
```bash
./bot
```

Ensure that your bot got neccesary permission for sending messages.
