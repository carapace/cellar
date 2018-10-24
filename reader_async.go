package cellar

import (
	"context"
)

type Rec struct {
	Data     []byte
	ChunkPos int64
	StartPos int64
	NextPos  int64
}

// ScanAsync runs Reader.Scan in a goroutine, returning the values obtained.
//
// ScanAsync honors context cancellations. If an error is received in the error channel, no more values will
// be scanned and the routine exits.
func (reader *Reader) ScanAsync(ctx context.Context, buffer int) (chan Rec, chan error) {
	vals := make(chan Rec, buffer)
	errs := make(chan error)

	go func() {
		// make sure we terminate the channel on scan read
		defer close(vals)
		defer close(errs)

		err := reader.Scan(func(ri *ReaderInfo, data []byte) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				vals <- Rec{data, ri.ChunkPos, ri.StartPos, ri.NextPos}
				return nil
			}
		})

		if err != nil {
			errs <- err
		}
	}()
	return vals, errs
}
