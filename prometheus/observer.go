package prometheus

// Observer is the interface that wraps the Observe method, which is used by
// Histogram and Summary to add observations.
type Observer interface {
	Observe(float64)
}

// The ObserverFunc type is an adapter to allow the use of ordinary
// functions as Observers. If f is a function with the appropriate
// signature, ObserverFunc(f) is an Observer that calls f.
//
// This adapter is usually used in connection with the Timer type, and there are
// two general use cases:
//
// The most common one is to use a Gauge as the Observer for a Timer.
// See the "Gauge" Timer example.
//
// The more advanced use case is to create a function that dynamically decides
// which Observer to use for observing the duration. See the "Complex" Timer
// example.
type ObserverFunc func(float64)

// Observe calls f(value). It implements Observer.
func (f ObserverFunc) Observe(value float64) {
	f(value)
}

// ObserverVec implements MetricVec for the purpose of using instance labels
// for an Observer.
type ObserverVec struct {
	*MetricVec
	Observer
}

// GetMetricWithLabelValues replaces the method of the same name in
// MetricVec. The difference is that this method returns a Summary and not a
// Metric so that no type conversion is required.
func (m *ObserverVec) GetMetricWithLabelValues(lvs ...string) (Observer, error) {
	metric, err := m.MetricVec.GetMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Observer), err
	}
	return nil, err
}

// GetMetricWith replaces the method of the same name in MetricVec. The
// difference is that this method returns a Observer and not a Metric so that no
// type conversion is required.
func (m *ObserverVec) GetMetricWith(labels Labels) (Observer, error) {
	metric, err := m.MetricVec.GetMetricWith(labels)
	if metric != nil {
		return metric.(Observer), err
	}
	return nil, err
}

// WithLabelValues works as GetMetricWithLabelValues, but panics where
// GetMetricWithLabelValues would have returned an error. By not returning an
// error, WithLabelValues allows shortcuts like
//     myVec.WithLabelValues("404", "GET").Observe(42.21)
func (m *ObserverVec) WithLabelValues(lvs ...string) Observer {
	return m.MetricVec.WithLabelValues(lvs...).(Observer)
}

// With works as GetMetricWith, but panics where GetMetricWithLabels would have
// returned an error. By not returning an error, With allows shortcuts like
//     myVec.With(Labels{"code": "404", "method": "GET"}).Observe(42.21)
func (m *ObserverVec) With(labels Labels) Observer {
	return m.MetricVec.With(labels).(Observer)
}