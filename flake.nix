{
  description = "A simple chat application for LAN";

  # Nixpkgs / NixOS version to use.
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, utils }: utils.lib.eachDefaultSystem (system:
    let
      version = "0.1.0";
      pname = "lan-chat";
      
      # Go in pkgs was 1.16
      goOverlay = final: prev: {
        go = prev.go.overrideAttrs rec { 
          version = "1.23.2";

          src = pkgs.fetchurl {
            url = "https://go.dev/dl/go${version}.src.tar.gz";
            hash = "sha256-NpMBYqk99BfZC9IsbhTa/0cFuqwrAkGO3aZxzfqc0H8=";
          };
        }; 
      };
      
      pkgs = import nixpkgs { inherit system; overlays = [ goOverlay ]; };
    in
    {
      packages = rec {
        lan-chat = pkgs.buildGoModule {
          inherit pname version;
          src = ./.;
          vendorHash = null;
        };
        default = lan-chat;
      };

      devShells.default = with pkgs; pkgs.mkShell {
        buildInputs = [ go ];
      };
    });
}
