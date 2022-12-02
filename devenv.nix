{ pkgs, ... }:

{
  # https://devenv.sh/packages/
  packages = [
    pkgs.git
    pkgs.go
    pkgs.protobuf
    pkgs.skaffold
    pkgs.kind
    pkgs.kubectl
    pkgs.kustomize
    pkgs.python310
    pkgs.poetry
  ];

}
