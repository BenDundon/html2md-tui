{
  description = "A TUI wrapper for html-to-markdown";

  inputs = {
    # Nixpkgs provides the Nix packages collection.
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";

    # Flake-utils helps with easier flake development.
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: flake-utils.lib.eachDefaultSystem (system: let
    pkgs = import nixpkgs { system = system; };

    goModule = pkgs.buildGoModule {
    pname = "html2md-tui";
    version = "0.5.0";

    src = pkgs.fetchFromGitHub {
      owner = "BenDundon";
      repo = "html2md-tui";
      rev = "v0.5.0";
      hash = nixpkgs.lib.fakeHash;
    };

    vendorHash = nixpkgs.lib.fakeHash;

    meta = {
      description = "A TUI Wrapper for html-to-markdown";
      homepage = "https://github.com/BenDundon/html2md-tui";
      license = nixpkgs.lib.licenses.gpl3;
    };
  };
  in {
    packages = {
      default = goModule;
    };
  });
}
