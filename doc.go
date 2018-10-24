/*
package cellar implements a low level, append only data storage engine.

cellar was originally forked from abdullin/cellar. This version focusses on a more standard API, less built in features
and more configurability. It is used in carapace/core/pkg/append-db (versioned object DB) as the storage engine.

It is completely embedded, internally buffering writes until chunks are filled up, then automatically flushing them to
file. See the example for a full usage of cellar.
*/
package cellar
