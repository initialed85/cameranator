package utils

import (
	"os"
	"os/signal"
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/relvacode/iso8601"
)

func GetUUID() uuid.UUID {
	someUUID, _ := uuid.NewRandom()

	return someUUID
}

func GetISO8601Time(rawTimestamp string) iso8601.Time {
	timestamp, _ := iso8601.ParseString(rawTimestamp)

	return iso8601.Time{Time: timestamp}
}

func WaitForCtrlC() {
	var wg sync.WaitGroup

	wg.Add(1)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		wg.Done()
	}()

	wg.Wait()
}

func Dereference(thing interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(thing)).Interface()
}
