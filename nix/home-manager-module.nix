self:
{ config, lib, pkgs, ... }:

let
  cfg = config.programs.goatty;
  yaml = pkgs.formats.yaml { };
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

    # Go's yaml decoder lower-cases struct field names, so the keys here are
    # all lower case with no separator: `brightblack`, not `brightBlack`.
    settings = lib.mkOption {
      type = yaml.type;
      default = { };
      example = lib.literalExpression ''
        {
          opacity = 0.9;
          font = {
            family = "JetBrains Mono";
            size = 14.0;
            ligatures = true;
          };
        }
      '';
      description = ''
        Written to {file}`$XDG_CONFIG_HOME/goatty/config.yaml`.
        Left alone when empty, so goatty falls back to its own defaults.
      '';
    };

    theme = lib.mkOption {
      type = yaml.type;
      default = { };
      example = lib.literalExpression ''
        {
          background = "#1d1f21";
          foreground = "#c5c8c6";
          brightblack = "#666666";
        }
      '';
      description = ''
        Written to {file}`$XDG_CONFIG_HOME/goatty/theme.yaml`.
        Left alone when empty, so goatty falls back to its built-in theme.
      '';
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [ cfg.package ];

    xdg.configFile."goatty/config.yaml" = lib.mkIf (cfg.settings != { }) {
      source = yaml.generate "goatty-config.yaml" cfg.settings;
    };

    xdg.configFile."goatty/theme.yaml" = lib.mkIf (cfg.theme != { }) {
      source = yaml.generate "goatty-theme.yaml" cfg.theme;
    };
  };
}
