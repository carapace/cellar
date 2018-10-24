# Cellar

[![Build Status](https://travis-ci.com/carapace/cellar.svg?branch=master)](https://travis-ci.com/carapace/cellar)
[![CircleCI](https://circleci.com/gh/carapace/cellar/tree/master.svg?style=svg)](https://circleci.com/gh/carapace/cellar/tree/master)
[![Coverage Status](https://coveralls.io/repos/github/carapace/cellar/badge.svg?branch=master)](https://coveralls.io/github/carapace/cellar?branch=master)

[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![](https://godoc.org/github.com/carapace/cellar?status.svg)](http://godoc.org/github.com/carapace/cellar)
[![Go Report Card](https://goreportcard.com/badge/github.com/carapace/cellar)](https://goreportcard.com/report/github.com/carapace/cellar)

Cellar is the append-only storage backend in Go designed based on Abdullin Cellar.
This fork is currently being redesigned, so the API should be considered unstable.

Core features:

- events are automatically split into the chunks;
- chunks may be encrypted using the Cipher interface;
- designed for batching operations (high throughput);
- supports single writer and multiple concurrent readers;
- store secondary indexes, lookups in the metadata DB.

# Contributors

In the alphabetical order:

- [Karel L. Kubat](https://github.com/KaiserKarel)
- [Rinat Abdullin](https://github.com/abdullin)

Don't hesitate to send a PR to include your profile.

# Design

Cellar stores data in a very simple manner:

- MetaDB database is used for keeping metadata (including user-defined), see metadb.go;
- a single pre-allocated file is used to buffer all writes;
- when buffer fills, it is compressed, encrypted and added to the chunk list.

# License

3-clause BSD license.
