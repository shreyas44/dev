{ pkgs ? import <nixpkgs> { } }:

pkgs.buildGoModule {
  name = "dev";
  src = pkgs.lib.cleanSource ./.;
  vendorSha256 = "jUvQCKiH/8336QXU45QZu1xFymilXHCH4O+jk8tLDWQ=";
}

