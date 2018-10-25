# Changelog
All notable changes to this project will be documented in this file.

## [Upcomming]
- automatic cleanup of lz4 buffer files
- integrity checks on startup
- easier-to-use indexing abstractions
- LMDB implementation of MetaDB

## [0.0.1] - 2018-10-25
### Added
- Cipher interface
- Compressor interface
- MetaDB interface
- BoltDB implementation of MetaDB
- Filelock
- God level DB struct with options
- Pluggable zap logger
- Async reader implements context cancellation
- extended test coverage

### Changed
- removed hardcoded encryption and compression
- removed hardcoded LMDB as metadb
- minor breaking API changes with abdullin/cellar

## [0.0.0] - abdullin/cellar
original fork of abdullin/cellar