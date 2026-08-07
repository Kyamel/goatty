{ lib
, buildGoModule
, pkg-config
, makeWrapper
, makeDesktopItem
, copyDesktopItems
, libGL
, libx11
, libxcursor
, libxi
, libxinerama
, libxrandr
, libxxf86vm
, src
, version
}:

let
  graphicsLibs = [
    libGL
    libx11
    libxcursor
    libxi
    libxinerama
    libxrandr
    libxxf86vm
  ];
in
buildGoModule {
  pname = "darktile";
  inherit src version;

  # The dependencies are committed under vendor/.
  vendorHash = null;

  subPackages = [ "cmd/darktile" ];

  nativeBuildInputs = [ pkg-config makeWrapper copyDesktopItems ];
  buildInputs = graphicsLibs;

  ldflags = [
    "-s"
    "-w"
    "-X github.com/liamg/darktile/internal/app/darktile/version.Version=${version}"
  ];

  # ebiten resolves libGL with dlopen rather than linking it, so nothing here
  # ends up in the binary's DT_NEEDED and the loader has no way to find it.
  postInstall = ''
    wrapProgram $out/bin/darktile \
      --prefix LD_LIBRARY_PATH : ${lib.makeLibraryPath graphicsLibs}
  '';

  desktopItems = [
    (makeDesktopItem {
      name = "darktile";
      exec = "darktile";
      icon = "utilities-terminal";
      desktopName = "Darktile";
      comment = "GPU rendered terminal emulator";
      categories = [ "System" "TerminalEmulator" ];
      terminal = false;
    })
  ];

  meta = {
    description = "GPU rendered terminal emulator";
    homepage = "https://github.com/liamg/darktile";
    license = lib.licenses.mit;
    mainProgram = "darktile";
    platforms = lib.platforms.linux;
  };
}
