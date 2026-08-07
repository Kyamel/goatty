{
  description = "darktile - GPU rendered terminal emulator";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go
            pkg-config
            git
            gopls
            delve
            goreleaser

            # conformance and performance test suites
            python3
            goperf # provides benchstat
            graphviz # pprof renders call graphs through dot
          ];

          buildInputs = with pkgs; [
            libGL
            libx11
            libxcursor
            libxi
            libxinerama
            libxrandr
            libxxf86vm
          ];
        };
      });
    };
}
