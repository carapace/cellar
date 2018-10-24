package app

import (
	"bytes"
	"context"
	"errors"
	"github.com/carapace/cellar"
	"runtime"
	"time"
)

func Routine(ctx context.Context, duration time.Duration) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()

	db, err := NewDB()
	if err != nil {
		return err
	}

	defer db.Close()

	now := time.Now()
	msgChan := make(chan []byte)

	// start a reader routine
	go func() {
		for {
			runtime.Gosched()
			reader := db.Reader() // reader needs to be called to get the new writes into the buffer
			msg := <-msgChan
			var found bool
			reader.Scan(func(pos *cellar.ReaderInfo, data []byte) error {
				if bytes.Equal(data, msg) {
					found = true
				}
				return nil
			})
			if !found {
				panic("unable to find msg")
			}
		}
	}()

	for {
		if time.Now().After(now.Add(duration)) {
			return nil // happy exit
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			runtime.Gosched()
			msg := []byte(time.Now().String())
			_, err := db.Append(msg)
			if err != nil {
				return err
			}
			db.Flush()
			msgChan <- msg
		}
	}
}
