# Changelog

All notable changes to this project will be documented in this file.

## [0.1.6](https://github.com/schmoli/cli-tools/compare/v0.1.5...v0.1.6) (2025-12-05)


### Bug Fixes

* add trans-cli to CI and release workflows ([#11](https://github.com/schmoli/cli-tools/issues/11)) ([bbd80bb](https://github.com/schmoli/cli-tools/commit/bbd80bb61607d159d88ea653b7239fb59e9c7bfd))

## [0.1.5](https://github.com/schmoli/cli-tools/compare/v0.1.4...v0.1.5) (2025-12-05)


### Features

* add trans-cli for Transmission RPC ([#9](https://github.com/schmoli/cli-tools/issues/9)) ([b99561a](https://github.com/schmoli/cli-tools/commit/b99561af8bba85c18c32051a8b832385d138ffb9))

## [0.1.4](https://github.com/schmoli/cli-tools/compare/v0.1.3...v0.1.4) (2025-12-04)


### Features

* build downloadable artifacts on main push for pre-release testing ([4a1c973](https://github.com/schmoli/cli-tools/commit/4a1c97356d5c7d93958b05fd86396499f77a1404))

## [0.1.3](https://github.com/schmoli/cli-tools/compare/v0.1.2...v0.1.3) (2025-12-04)


### Bug Fixes

* build binaries in release-please workflow when release created ([b7aa7a0](https://github.com/schmoli/cli-tools/commit/b7aa7a0a6bf5898d1b4b22c1820df3208e948ad5))

## [0.1.2](https://github.com/schmoli/cli-tools/compare/v0.1.1...v0.1.2) (2025-12-04)


### Features

* add --insecure/-k flag to skip TLS verification ([c441911](https://github.com/schmoli/cli-tools/commit/c44191165129a5eff7d35ec2e584477947da24e0))
* add install test script for Docker validation ([c41034e](https://github.com/schmoli/cli-tools/commit/c41034eecd05ed7e8ca541a0390884e4056b0724))
* add nproxy-cli for nginx-proxy-manager API ([86016c8](https://github.com/schmoli/cli-tools/commit/86016c8cd257a01037c2fd75b6a5866dcef3e685))
* add release pipeline, install script, version embedding ([0f4cf14](https://github.com/schmoli/cli-tools/commit/0f4cf14a9df58e34292fbe82861810b0675363ef))
* **go:** add API and output models with type mappings ([239f157](https://github.com/schmoli/cli-tools/commit/239f157e3cf960531a725f0f1bdd14b679c3cb7a))
* **go:** add error types with exit codes ([b1f1774](https://github.com/schmoli/cli-tools/commit/b1f1774d65e0287486ee98ce529974a162a83011))
* **go:** add HTTP client with auth and error handling ([63c9b53](https://github.com/schmoli/cli-tools/commit/63c9b53eb5a38b93b09d1cb359485f48bd288518))
* **go:** add YAML output formatting ([fbe8891](https://github.com/schmoli/cli-tools/commit/fbe8891e903129a8628267ca668faf2f9114d6a2))
* **go:** implement CLI with stacks and endpoints commands ([f628c79](https://github.com/schmoli/cli-tools/commit/f628c79ff1e3390ffb9fb3700d152ef802b241d9))
* **go:** init module with cmd and pkg structure ([708196b](https://github.com/schmoli/cli-tools/commit/708196b1c9a4e311685c5fa12c7edc3d55f1f625))
* **rust:** add API and output models with type mappings ([8482eed](https://github.com/schmoli/cli-tools/commit/8482eed9e259d6e423c60f240a82435944a780f3))
* **rust:** add error types with exit codes ([5757e4d](https://github.com/schmoli/cli-tools/commit/5757e4dc0d46247b22e086a7e043d7147ea00f2e))
* **rust:** add HTTP client with auth and error handling ([cfc1d72](https://github.com/schmoli/cli-tools/commit/cfc1d72751dd014b51c6d49a09846a0ece494e2f))
* **rust:** add YAML output formatting ([78f007f](https://github.com/schmoli/cli-tools/commit/78f007fc9e122251ebe8e39312712ca30415e06d))
* **rust:** implement CLI with stacks and endpoints commands ([481c480](https://github.com/schmoli/cli-tools/commit/481c4801fa0c5f9039a9830c7783ba50673ac324))
* **rust:** init workspace with lib and cli crates ([bbebd50](https://github.com/schmoli/cli-tools/commit/bbebd5037c365b29d9251b749b588bcf4dd6df0d))


### Bug Fixes

* add workflow_dispatch to release-please ([0e8b9b3](https://github.com/schmoli/cli-tools/commit/0e8b9b39fd741ae762f0653a505c0aaf9c9fa973))
* code quality improvements, enhanced build script ([a28b7d4](https://github.com/schmoli/cli-tools/commit/a28b7d4e53bb2e1ba6eb58b0e0b178f1bdd60d12))
* endpoint status uses active/inactive per requirements ([4388cdf](https://github.com/schmoli/cli-tools/commit/4388cdf6388af3c8c84743e47c5f225b1f74a9ab))
* restore --insecure flag, add unit tests ([961d8b7](https://github.com/schmoli/cli-tools/commit/961d8b7df48f73a5f0b4c78e6c6cdce5eb42ad29))
* set initial version to 0.1.0 ([05c9ece](https://github.com/schmoli/cli-tools/commit/05c9ece0489d3aef9344cdb386f9b0c9191ecb73))

## [0.1.1](https://github.com/schmoli/cli-tools/compare/v0.1.0...v0.1.1) (2025-12-04)


### Features

* add --insecure/-k flag to skip TLS verification ([c441911](https://github.com/schmoli/cli-tools/commit/c44191165129a5eff7d35ec2e584477947da24e0))
* add nproxy-cli for nginx-proxy-manager API ([86016c8](https://github.com/schmoli/cli-tools/commit/86016c8cd257a01037c2fd75b6a5866dcef3e685))
* add release pipeline, install script, version embedding ([0f4cf14](https://github.com/schmoli/cli-tools/commit/0f4cf14a9df58e34292fbe82861810b0675363ef))
* **go:** add API and output models with type mappings ([239f157](https://github.com/schmoli/cli-tools/commit/239f157e3cf960531a725f0f1bdd14b679c3cb7a))
* **go:** add error types with exit codes ([b1f1774](https://github.com/schmoli/cli-tools/commit/b1f1774d65e0287486ee98ce529974a162a83011))
* **go:** add HTTP client with auth and error handling ([63c9b53](https://github.com/schmoli/cli-tools/commit/63c9b53eb5a38b93b09d1cb359485f48bd288518))
* **go:** add YAML output formatting ([fbe8891](https://github.com/schmoli/cli-tools/commit/fbe8891e903129a8628267ca668faf2f9114d6a2))
* **go:** implement CLI with stacks and endpoints commands ([f628c79](https://github.com/schmoli/cli-tools/commit/f628c79ff1e3390ffb9fb3700d152ef802b241d9))
* **go:** init module with cmd and pkg structure ([708196b](https://github.com/schmoli/cli-tools/commit/708196b1c9a4e311685c5fa12c7edc3d55f1f625))
* **rust:** add API and output models with type mappings ([8482eed](https://github.com/schmoli/cli-tools/commit/8482eed9e259d6e423c60f240a82435944a780f3))
* **rust:** add error types with exit codes ([5757e4d](https://github.com/schmoli/cli-tools/commit/5757e4dc0d46247b22e086a7e043d7147ea00f2e))
* **rust:** add HTTP client with auth and error handling ([cfc1d72](https://github.com/schmoli/cli-tools/commit/cfc1d72751dd014b51c6d49a09846a0ece494e2f))
* **rust:** add YAML output formatting ([78f007f](https://github.com/schmoli/cli-tools/commit/78f007fc9e122251ebe8e39312712ca30415e06d))
* **rust:** implement CLI with stacks and endpoints commands ([481c480](https://github.com/schmoli/cli-tools/commit/481c4801fa0c5f9039a9830c7783ba50673ac324))
* **rust:** init workspace with lib and cli crates ([bbebd50](https://github.com/schmoli/cli-tools/commit/bbebd5037c365b29d9251b749b588bcf4dd6df0d))


### Bug Fixes

* add workflow_dispatch to release-please ([0e8b9b3](https://github.com/schmoli/cli-tools/commit/0e8b9b39fd741ae762f0653a505c0aaf9c9fa973))
* code quality improvements, enhanced build script ([a28b7d4](https://github.com/schmoli/cli-tools/commit/a28b7d4e53bb2e1ba6eb58b0e0b178f1bdd60d12))
* endpoint status uses active/inactive per requirements ([4388cdf](https://github.com/schmoli/cli-tools/commit/4388cdf6388af3c8c84743e47c5f225b1f74a9ab))
* restore --insecure flag, add unit tests ([961d8b7](https://github.com/schmoli/cli-tools/commit/961d8b7df48f73a5f0b4c78e6c6cdce5eb42ad29))
* set initial version to 0.1.0 ([05c9ece](https://github.com/schmoli/cli-tools/commit/05c9ece0489d3aef9344cdb386f9b0c9191ecb73))

## [Unreleased]

### Added
- Initial release of cli-tools monorepo
- `portainer-cli` - CLI for Portainer API
- `nproxy-cli` - CLI for nginx-proxy-manager API
- Automated install script
- Cross-platform builds (macOS, Linux)
