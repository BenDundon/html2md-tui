# Installation

This project is inteded for use with the [Nix Package Manager](nixos.org), and was developed on NixOS.

## Running from source

To run locally, ensure you have enabled [flakes](https://wiki.nixos.org/wiki/Flakes), then pull the repo and run

```
nix run  
```
in the project root directory.

### Without flakes

If you don't like flakes, you can run

```
  nix-build
```
in the project root directory, and then run the binary from `/result/bin/html2md-tui`

## Installing on NixOS

(Note: Currently this doesn't work)

To install to `$PATH` on NixOS, copy `html2md-tui.nix` to `/etc/nixos/flakes/html2md-tui/flake.nix`. Then, add the following to your `configuration.nix`:

```nix
let
  html2md-tui = builtins.getFlake "/etc/nixos/flakes/html2md-tui";
in
{
  # ...rest of your configuration.nix
}
  
```

and add `html2md-tui.packages.${pkgs.system}.default` to either your `users.<user>.packages` or `environment.systemPackages`.

## Installing on systems without Nix

For systems without nix, the program will have to be built from source. For this, you will need [Go](https://go.dev/) installed on your system. Then, run
```
go build  
```
to build the project. The binary will build in the project root, and may be moved to your `/usr/bin`, `.local/bin`, etc., as long as it's put somewhere on your $PATH.

## Installing for Windows

Please don't.
