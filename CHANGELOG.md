# Changelog

## [1.2.2](https://github.com/nathan-nicholson/note/compare/v1.2.1...v1.2.2) (2025-12-02)


### Bug Fixes

* add cross-compilation environment variables for ARM64 ([77c0d19](https://github.com/nathan-nicholson/note/commit/77c0d1931e5091434d6692152b226f4e31f631ae))
* add cross-compilation tools and temporarily build Linux-only ([0c4f7da](https://github.com/nathan-nicholson/note/commit/0c4f7da3738c7b9f4d7b3e778a2b3db2a20105e3))
* change CI to validate GoReleaser config instead of building ([a3e3ab5](https://github.com/nathan-nicholson/note/commit/a3e3ab5e184d8898bf49562c9219035c7fe6eb23))
* install cross-compilation tools in CI for binary builds ([fe7bf47](https://github.com/nathan-nicholson/note/commit/fe7bf471a7a9951ac0fbb7f7d7b7875e5483fd21))
* restore workflows to zig-based cross-compilation ([97af8f1](https://github.com/nathan-nicholson/note/commit/97af8f1f40339ee9c44cf9b3c516dddae061c9d8))
* simplify CI to build only Linux amd64 with CGO ([cc10a57](https://github.com/nathan-nicholson/note/commit/cc10a5779042f484a5bc35cf7426c5e3d080c4bf))
* simplify cross-compilation by disabling CGO ([61dc20f](https://github.com/nathan-nicholson/note/commit/61dc20f483dd9bc9a12ac4e40f0c72953dc727f6))
* sort notes in chronological order (oldest first) ([#16](https://github.com/nathan-nicholson/note/issues/16)) ([377c4cf](https://github.com/nathan-nicholson/note/commit/377c4cf7c8dc123869b57714f17c3bf4143348e8))
* use native runners for macOS and Linux builds ([9c97a3c](https://github.com/nathan-nicholson/note/commit/9c97a3c76eded970b27485ddf51eef5c419a846a))
* use zig for CGO cross-compilation to support macOS ([1b459b0](https://github.com/nathan-nicholson/note/commit/1b459b00b566fa2a3d67d27a1268464d8a6630a9))

## [1.2.1](https://github.com/nathan-nicholson/note/compare/v1.2.0...v1.2.1) (2025-12-01)


### Bug Fixes

* enable CGO for SQLite support and improve version handling ([feb51bf](https://github.com/nathan-nicholson/note/commit/feb51bf5e95a55bdbdb5643d019207353fae9738))

## [1.2.0](https://github.com/nathan-nicholson/note/compare/v1.1.3...v1.2.0) (2025-11-28)


### Features

* initial release of note CLI ([d092468](https://github.com/nathan-nicholson/note/commit/d09246806367ab84a107e05d9750f3a27bb0a407))


### Bug Fixes

* add local timezone support to features list and fix gorelease flow ([#6](https://github.com/nathan-nicholson/note/issues/6)) ([27d4337](https://github.com/nathan-nicholson/note/commit/27d4337216f1ba55a0a5935012d42a445e581f89))
* configure GoReleaser archives to include binaries ([#10](https://github.com/nathan-nicholson/note/issues/10)) ([83f3dbe](https://github.com/nathan-nicholson/note/commit/83f3dbe3755fb3e21681eba9460a95347e4c7aba))
* correct timestamp timezone handling ([#4](https://github.com/nathan-nicholson/note/issues/4)) ([bd6ce29](https://github.com/nathan-nicholson/note/commit/bd6ce29222ff33a88f19668c597c838ef0fbce88))
* GoReleaser config, release workflow, and end-of-week edge case ([#12](https://github.com/nathan-nicholson/note/issues/12)) ([9951327](https://github.com/nathan-nicholson/note/commit/995132757fb877b45313f93359569bf4d305f568))
* integrate GoReleaser into release-please workflow ([05e3cb9](https://github.com/nathan-nicholson/note/commit/05e3cb9a807cd447c54a11f8e6d6f49f06cb86f7))
* update goreleaser config and add dist to gitignore ([92dac98](https://github.com/nathan-nicholson/note/commit/92dac980c5974fa0fff34fc7bc11a6af18ec82f6))
* update goreleaser config for v2 compatibility ([b0eb151](https://github.com/nathan-nicholson/note/commit/b0eb151f780157b81434ddfa2d5c300c2a26d37b))
* update goreleaser workflow trigger to published event ([0fa0faf](https://github.com/nathan-nicholson/note/commit/0fa0faf5bddcd9c5fc2185c10c17b801b7e4f2c3))
* use default GITHUB_TOKEN for release-please ([#5](https://github.com/nathan-nicholson/note/issues/5)) ([f8f5d8d](https://github.com/nathan-nicholson/note/commit/f8f5d8d6c8d6c563c6dfc7ba5f6a3d16ceaa6ffc))
* use PAT for release-please to trigger release workflow ([4f7b9a4](https://github.com/nathan-nicholson/note/commit/4f7b9a457aa7c651b637e0fc3734513daafd7ea9))

## [1.1.3](https://github.com/nathan-nicholson/note/compare/v1.1.2...v1.1.3) (2025-11-28)


### Bug Fixes

* update goreleaser workflow trigger to published event ([0fa0faf](https://github.com/nathan-nicholson/note/commit/0fa0faf5bddcd9c5fc2185c10c17b801b7e4f2c3))

## [1.1.2](https://github.com/nathan-nicholson/note/compare/v1.1.1...v1.1.2) (2025-11-28)


### Bug Fixes

* GoReleaser config, release workflow, and end-of-week edge case ([#12](https://github.com/nathan-nicholson/note/issues/12)) ([9951327](https://github.com/nathan-nicholson/note/commit/995132757fb877b45313f93359569bf4d305f568))

## [1.1.1](https://github.com/nathan-nicholson/note/compare/v1.1.0...v1.1.1) (2025-11-27)


### Bug Fixes

* configure GoReleaser archives to include binaries ([#10](https://github.com/nathan-nicholson/note/issues/10)) ([83f3dbe](https://github.com/nathan-nicholson/note/commit/83f3dbe3755fb3e21681eba9460a95347e4c7aba))

## [1.1.0](https://github.com/nathan-nicholson/note/compare/v1.0.2...v1.1.0) (2025-11-27)


### Features

* initial release of note CLI ([d092468](https://github.com/nathan-nicholson/note/commit/d09246806367ab84a107e05d9750f3a27bb0a407))


### Bug Fixes

* add local timezone support to features list and fix gorelease flow ([#6](https://github.com/nathan-nicholson/note/issues/6)) ([27d4337](https://github.com/nathan-nicholson/note/commit/27d4337216f1ba55a0a5935012d42a445e581f89))
* correct timestamp timezone handling ([#4](https://github.com/nathan-nicholson/note/issues/4)) ([bd6ce29](https://github.com/nathan-nicholson/note/commit/bd6ce29222ff33a88f19668c597c838ef0fbce88))
* update goreleaser config and add dist to gitignore ([92dac98](https://github.com/nathan-nicholson/note/commit/92dac980c5974fa0fff34fc7bc11a6af18ec82f6))
* update goreleaser config for v2 compatibility ([b0eb151](https://github.com/nathan-nicholson/note/commit/b0eb151f780157b81434ddfa2d5c300c2a26d37b))
* use default GITHUB_TOKEN for release-please ([#5](https://github.com/nathan-nicholson/note/issues/5)) ([f8f5d8d](https://github.com/nathan-nicholson/note/commit/f8f5d8d6c8d6c563c6dfc7ba5f6a3d16ceaa6ffc))
* use PAT for release-please to trigger release workflow ([4f7b9a4](https://github.com/nathan-nicholson/note/commit/4f7b9a457aa7c651b637e0fc3734513daafd7ea9))

## Changelog
