### Fixed

- SSH `loginShell` starts work when the account shell is fish as well as a POSIX shell. <!-- fix: e3d7d18f2 -->
- An X11 `DISPLAY` with an explicit `unix:` host is rejected as `xcb_connect` rejects it; `unix/:` remains the local transport. <!-- fix: e3d7d18f2 -->
- Historical `Eire` `Date` strings use the names the qualified Node prints. <!-- fix: e3d7d18f2 -->
