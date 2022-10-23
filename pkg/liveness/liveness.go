package liveness

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type HasLiveness interface {
	IsLive() bool
}

type Agent struct {
	services []HasLiveness
	serveMux *http.ServeMux
	server   *http.Server
}

func Open(
	services []HasLiveness,
	port int64,
) (*Agent, error) {
	a := Agent{
		services: services,
		serveMux: http.NewServeMux(),
		server:   nil,
	}

	a.serveMux.HandleFunc(
		"/healthz",
		a.handle,
	)

	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: a.serveMux,
	}

	errors := make(chan error)

	t := time.NewTimer(time.Millisecond * 100)

	go func() {
		go func() {
			<-t.C
			errors <- nil
		}()

		err := a.server.ListenAndServe()
		if err != nil {
			errors <- err
		}

		t.Stop()
	}()

	runtime.Gosched()

	time.Sleep(time.Millisecond * 100)

	err := <-errors
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *Agent) handle(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, server := range a.services {
		isLive := server.IsLive()
		if !isLive {
			responseWriter.WriteHeader(http.StatusBadGateway)
			return
		}
	}

	responseWriter.WriteHeader(http.StatusOK)
}

func (a *Agent) Close() {
	_ = a.server.Close()
}
