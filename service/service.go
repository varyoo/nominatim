package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/varyoo/nominatim"
)

type Job interface {
	Search() url.Values
	SetCoordinates(jsonResponse io.ReadCloser) error
}

type Service struct {
	Nominatim nominatim.Client

	jobs chan Job
	stop chan bool
}

func New(client nominatim.Client) Service {
	s := Service{
		Nominatim: client,

		// 3 jobs max in queue
		jobs: make(chan Job, 3),

		stop: make(chan bool),
	}

	go s.Go()

	return s
}

func (s Service) work(job Job) error {
	search := job.Search()

	body, err := s.Nominatim.Lookup(search)
	if err != nil {
		return err
	}

	if err = job.SetCoordinates(body); err != nil {
		return fmt.Errorf("set coordinates: %s", err.Error())
	}

	return nil
}

func (s Service) Go() {
	for {
		select {
		case job := <-s.jobs:
			if err := s.work(job); err != nil {
				log.Println("Nominatim service:", err)
			}

			time.Sleep(time.Second * 5)

		case <-s.stop:
			close(s.stop)
			return
		}
	}
}

func (s Service) Localize(ctx context.Context, job Job) error {
	select {
	case s.jobs <- job:
		return nil

	case <-ctx.Done():
		// job dropped
		return ctx.Err()
	}
}

func (s Service) Close() {
	s.stop <- true

	// wait for current job to finish
	<-s.stop
}
