# Dotty

A dead simple dotfiles manager.

Use at your own risk.

## install

```sh
brew tap codethread/homebrew-dotty
brew install dotty
```

## usage

```sh
dotty help
```

## contributing

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
