### Changed

- Fork releases are published only after the Linux arm64, macOS arm64 and Windows amd64 archives run natively (`pig --version` and `pig update`). Linux amd64 still runs on the build host; the remaining targets are cross-built and the release notes say so.
