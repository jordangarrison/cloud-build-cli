let
  unstable = import
    (fetchTarball "https://nixos.org/channels/nixos-unstable/nixexprs.tar.xz")
    { };
in
{ nixpkgs ? import <nixpkgs> { } }:
with nixpkgs;
mkShell {
  name = "cluster-health-service";

  buildInputs = with pkgs; [ unstable.go_1_18 unstable.docker ];
}
