# `stylish`

![Logo for stylish](./assets/logo.png)

Welcome to `stylish`, a simple and intuitive way to create `dircolors`-compatible config files! This affects programs like `ls`, `tree`, `fd`, and any other tools that opt to respect the `LS_COLORS` environment variable.

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

- Have `brew` installed
- Run the following:
```sh
brew tap daltonsw/packages
brew install stylish
```

### From Go

- Have `Go` 
- Have your `Go` install location on your `$PATH`
- Run the following: 
```sh
go install go.dalton.dog/stylish@latest
```

## Usage

- Start the program with `stylish`. This will:
    - Create a `stylish` directory in your user's default config directory (typically `~/.config`)
    - Create a `default` theme inside of that directory. **Note:** This theme is intended to be used on a dark background
- With the program running, you're able to create and edit your themes to your heart's content
- Once you're ready to apply a theme, you'll need to add the following to your shell's init file (`~/.bashrc`, `~/.zshrc`, etc.):
    - **Required:** `eval $(stylish apply <theme>)`
    - *Recommended:* `alias ls=ls --color=auto`
- Once your init file is edited, relaunch your shell to start seeing the updated colors.

### P.S.

Want to handle your hex code journey in your terminal too? Check out [termpicker](https://github.com/ChausseBenjamin/termpicker)!

## Shoutouts

- [Jess](https://jessicakasper.com) for the great banner!
- [Vivid](https://github.com/sharkdp/vivid) for being a great program and a great reference
- [Catppuccin](https://github.com/catppuccin) for having pretty palettes
- [CharmBracelet](https://github.com/charmbracelet) for the amazing modules for style, form, and function

## Contributions

Contributions are very welcome! I'd love for y'all to contribute themes you develop as well as expanding on the defaults to make them more reasonable. Check out the [CONTRIBUTING](./CONTRIBUTING.md) file for specifics.

## License

Copyright 2025 - Dalton Williams  
Check [LICENSE](./LICENSE.md) in repo for full details
