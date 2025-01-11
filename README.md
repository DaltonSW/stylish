# stylish

Welcome to `stylish`, a simple and intuitive way to create `dircolors`-compatible config files!  
This affects programs like `ls`, `tree`, `fd`, and any other tools that opt to respect the `LS_COLORS` environment variable.

## Why use `stylish`?

- NO file editing!
- NO dealing with encoding strings!
- NO manual mapping of hex codes to 8-bit colors!
- NO blindly working without an actual preview!
- NO need to hunt down scattered, poor documentation and references!

![Demo of stylish](./assets/demo.gif)

## Installation

### Github Releases

- Go to the `Releases` tab
- Download the latest binary for your OS
- Place it on your `$PATH` and ensure it is executable

### Homebrew

**Requirements:**
- Have `brew` installed

- Run the following:
```sh
brew tap daltonsw/packages
brew install stylish
```

### From Source

**Requirements:**
- Have `Go` installed
- Have your `Go` install location on your `$PATH`

- Clone the repo with `git clone https://github.com/DaltonSW/stylish.git`
- `cd` into the cloned directory
- Run `go mod tidy` to download module requirements
- Run `go install .` to install the `stylish` binary

## Usage

- The program is able to be run immediately after install. Start it with `stylish`. This will...
    - Create a `stylish` folder in your user's default config folder (typically `~/.config`)
    - Creates a `default` theme inside of that folder
- With the program running, you're able to create and edit your themes to your heart's content
- Once you're ready to apply a theme, you'll need to add the following to your shell's init file (`~/.bashrc`, `~/.zshrc`, etc.):
    - **Required:** `eval $(stylish apply <theme>)`
    - *Recommended:* `alias ls=ls --color=auto`

## P.S.

Want to handle your hex code journey in your terminal too? Check out [termpicker](https://github.com/ChausseBenjamin/termpicker)!
