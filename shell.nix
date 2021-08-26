let pkgs = import ./nix/nixpkgs.nix;
in pkgs.mkShell {
  packages = with pkgs; [ gnumake go terraform_1_0_0 goimports golangci-lint ];
}
