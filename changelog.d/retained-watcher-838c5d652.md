### Fixed

- A retained packed process now has one watcher for its lifetime instead of one per reload, and a crash is attributed to the current generation. 200 reloads leave one watched state rather than 201. <!-- fix: 838c5d652 -->
