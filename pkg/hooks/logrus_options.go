package hooks

import (
	"sort"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type logrusOptions []logrus.Level

type logrusOptionApplyFunc func(opts logrusOptions) logrusOptions

type LogrusOption interface {
	apply(opts logrusOptions) logrusOptions
}

func (f logrusOptionApplyFunc) apply(opts logrusOptions) logrusOptions {
	return f(opts)
}

func FromLogrusLevel(level logrus.Level) LogrusOption {
	levels := make([]logrus.Level, 0, len(logrus.AllLevels))
	add := false
	for i := range logrus.AllLevels {
		idx := len(logrus.AllLevels) - 1 - i
		l := logrus.AllLevels[idx]
		if add == false && l != level {
			continue
		}
		levels = append(levels, l)
		add = true
	}

	return WithLogrusLevels(levels)
}

func WithLogrusLevels(levels []logrus.Level) LogrusOption {
	return logrusOptionApplyFunc(func(opts logrusOptions) logrusOptions {

		for _, level := range levels {
			if slices.Contains(opts, level) == true {
				continue
			}
			opts = append(opts, level)
		}
		sort.Slice(opts, func(i, j int) bool {
			return opts[i] > opts[j]
		})

		return opts
	})
}

func newLogrusOptions(options ...LogrusOption) logrusOptions {
	var res logrusOptions
	for _, o := range options {
		res = o.apply(res)
	}
	return res
}
