package hooks

import (
	"testing"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func TestFromLogrusLevel(t *testing.T) {
	testdata := []struct {
		Level    logrus.Level
		Expected []logrus.Level
	}{
		{logrus.PanicLevel, []logrus.Level{logrus.PanicLevel}},
		{logrus.FatalLevel, []logrus.Level{logrus.FatalLevel, logrus.PanicLevel}},
		{logrus.ErrorLevel, []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}},
		{logrus.WarnLevel, []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}},
		{logrus.InfoLevel, []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}},
		{logrus.DebugLevel, []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}},
		{logrus.TraceLevel, []logrus.Level{logrus.TraceLevel, logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}},
	}

	for _, d := range testdata {
		opts := newLogrusOptions(FromLogrusLevel(d.Level))
		if len(opts) != len(d.Expected) {
			t.Errorf("mismatched size=%d for level %s. Expected: %d",
				len(opts), d.Level, len(d.Expected))
		}

		size := Min(len(d.Expected), len(opts))

		for i := 0; i < size; i++ {
			if opts[i] != d.Expected[i] {
				t.Errorf("mismatched level %s at %d for level %d, expected %s",
					opts[i], i, d.Level, d.Expected[i])
			}
		}

		for _, l := range opts[size:] {
			t.Errorf("unexpected level %s for %s", l, d.Level)
		}

		for _, l := range d.Expected[size:] {
			t.Errorf("missing level %s for %s", l, d.Level)
		}
	}

}
