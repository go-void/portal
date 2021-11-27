package server

import (
	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/collector"
	"github.com/go-void/portal/pkg/dio"
	"github.com/go-void/portal/pkg/filter"
	"github.com/go-void/portal/pkg/pack"
	"github.com/go-void/portal/pkg/resolver"
)

func WithFilter(f filter.Engine) OptionsFunc {
	return func(s *Server) error {
		s.Filter = f
		return nil
	}
}

func WithCache(c cache.Cache) OptionsFunc {
	return func(s *Server) error {
		s.cacheEnabled = true
		s.Cache = c
		return nil
	}
}

func WithResolver(r resolver.Resolver) OptionsFunc {
	return func(s *Server) error {
		s.Resolver = r
		return nil
	}
}

func WithCollector(c collector.Collector) OptionsFunc {
	return func(s *Server) error {
		s.Collector = c
		return nil
	}
}

func WithUnpacker(u pack.Unpacker) OptionsFunc {
	return func(s *Server) error {
		s.Unpacker = u
		return nil
	}
}

func WithPacker(p pack.Packer) OptionsFunc {
	return func(s *Server) error {
		s.Packer = p
		return nil
	}
}

func WithReader(r dio.Reader) OptionsFunc {
	return func(s *Server) error {
		s.Reader = r
		return nil
	}
}

func WithWriter(w dio.Writer) OptionsFunc {
	return func(s *Server) error {
		s.Writer = w
		return nil
	}
}

func WithAcceptFunc(a AcceptFunc) OptionsFunc {
	return func(s *Server) error {
		s.AcceptFunc = a
		return nil
	}
}
