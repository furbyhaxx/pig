### Changed

- Fork releases are published only after the Linux arm64 archive runs natively (`pig --version` and `pig update`); Linux amd64 runs on the build host and Android arm64 is ELF-checked.

### Removed

- macOS and Windows release archives. They are halted targets (SPEC-0001 D4); the installer refuses macOS and the update manifest lists only Linux and Android. Build from source on those systems.
