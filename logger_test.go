package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	testsWithoutErrors := []struct {
		desc     string
		loglevel string
		want     zapcore.Level
	}{
		{
			desc:     "Loglevel debug",
			loglevel: "debug",
			want:     zapcore.DebugLevel,
		},
		{
			desc:     "Loglevel info",
			loglevel: "info",
			want:     zapcore.InfoLevel,
		},
		{
			desc:     "Loglevel error",
			loglevel: "error",
			want:     zapcore.ErrorLevel,
		},
	}

	testsWithErrors := []struct {
		desc     string
		loglevel string
	}{
		{
			desc:     "Loglevel invalid",
			loglevel: "invalid",
		},
	}

	for _, test := range testsWithoutErrors {
		t.Run(test.desc, func(t *testing.T) {
			l, err := NewLogger(test.loglevel)
			assert.NoError(t, err)
			assert.True(t, l.Core().Enabled(test.want))
		})
	}

	for _, test := range testsWithErrors {
		t.Run(test.desc, func(t *testing.T) {
			_, err := NewLogger(test.loglevel)
			assert.Error(t, err)
		})
	}
}
