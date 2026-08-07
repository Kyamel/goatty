self:
{ config, lib, pkgs, ... }:

let
  cfg = config.programs.goatty;
in
{
  options.programs.goatty = {
    enable = lib.mkEnableOption "the goatty terminal emulator";

    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.goatty;
      defaultText = lib.literalExpression "goatty.packages.\${system}.goatty";
      description = "The goatty package to install.";
    };
  };

  # goatty only ever reads config.yaml and theme.yaml from the user's config
  # directory, so there is nothing for a system-wide module to configure. Use
  # the home-manager module for settings and themes.
  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];
  };
}
