{
  description = "darktile - GPU rendered terminal emulator";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      devShells = forAllSystems (pkgs:
        let
          graphicsLibs = with pkgs; [
            libGL
            libx11
            libxcursor
            libxi
            libxinerama
            libxrandr
            libxxf86vm
          ];
        in
        {
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
            xvfb-run # render goldens need a display
            imagemagick # and a way to diff the screenshots
            util-linux # flock, so two render runs cannot stomp each other
          ];

          buildInputs = graphicsLibs;

          # ebiten from v2.9 on resolves libGL with dlopen instead of linking
          # it, so having these as buildInputs is not enough to run anything.
          shellHook = ''
            export LD_LIBRARY_PATH="${pkgs.lib.makeLibraryPath graphicsLibs}''${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
          '';
        };
      });
    };
}
