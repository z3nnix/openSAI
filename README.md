

# OpenSAI <img src="Logo.png" width=40 style="position: absolute; bottom: 0">


Open-source implementation of SAI (SimpleAI) in Go.

![Static Badge](https://img.shields.io/badge/language-Go-blue) ![](https://img.shields.io/github/commit-activity/w/z3nnix/openSAI/main
) ![](https://img.shields.io/github/stars/z3nnix/openSAI?style=flat-square
)

# Setup

First, you need to install dependencies.

### In Debian-based (Ubuntu, etc.)
```bash
sudo apt install go
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
go build -o bot src/bot.go
```

Run it within
```bash
./bot
```

Ensure that your bot got neccesary permission for sending messages.
