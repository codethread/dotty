# Dotty

A dead simple dotfiles manager.

Use at your own risk.

## How it works

Here's how I manage my dotfiles, and if you do the same, then dotty could be for you!

I keep a [repo](https://github.com/codethread/PersonalConfigs) with all my dotfiles and `.config` folders, but it also has other bits and pieces to help get me up and running on any new machine. A subset of this looks like:

```
├── README.org
├── Session.vim
├── _colors+fonts
│   ├── iterm_tokyonight_storm.itermcolors
│   └── xterm-256color-italic.terminfo
├── .config
│   ├── kitty
│   │   ├── kitty.conf
│   │   └── open-actions.conf
│   └── nvim
│       ├── init.lua
│       └── lua
│           └── codethread
│               ├── keymaps.lua
│               └── lsp
│                   ├── init.lua
│                   └── settings
│                       ├── tailwindcss.lua
│                       └── tsserver.lua
├── .dottyignore
├── .gitconfig
├── .gitignore
├── .tmux.conf
├── .stylua.toml
├── .zshrc
├── Brewfile
├── README.org
└── _packages
    └── .piplist
```

Dotty's approach is simple, it just symlinks this entire tree, one file at a time into my HOME directory - creating folders if needs be, without disturbing any existing folders.

Not everything gets copied however, as I don't want everything cluttering my HOME, instead dotty will look for a `gitignore` and `gitignore_global` inside my dotfiles, and ignore any files listed there. It will also read from a `.dottyignore` file, which is a simple list of regexps, which can be used to omit any additional files. For example I have:

```sh
^.git$
^_.*
.gitignore$
README
^.stylua.toml$
^.dottyignore$
```

For me, a simple `_` prefix helps me keep track of what files and folders are not inteded to be linked. Not all files can follow that format, e.g `.stylua.toml`, hence `.dottyignore` allows more specific controls.

If a folder is marked as ignored, dotty will not traverse it

When dotty runs, it stores a list of created files in a temporary file, and then each time it runs, it removes these files. Empty folders will be deleted as it goes, meaning you are free to move and rename files/folders in your dotfolder as you please, and when dotty next runs, your HOME directory will be left in a clean state, before all the files are symlinked back in.

Personally I like to integrate this flow with my editor, see the wiki for an example; this means dotty is run everytime I save or delete a file in my editor (and dotty is extremely fast, so there's no issue of slowdown)

## Install

```sh
brew tap codethread/homebrew-dotty
brew install dotty
```

## Usage

```sh
dotty help
```

## Contributing

- run
  ```sh
  go run main.go <cmd> ...args
  # e.g.
  go run main.go setup --dry-run
  ```
- test
  ```sh
  go test ./...
  ```
- add command
  ```sh
  go install github.com/spf13/cobra-cli@latest
  cobra-cli add <cmd>
  ```

### Releasing

check the wiki
