package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"math"
	"testing"
)

func TestExtractFailureModel(t *testing.T) {
	// Known models must resolve without error and yield a usable draw function.
	for _, model := range []string{"weibull", "lognormal"} {
		draw, err := ExtractFailureModel(configextractor.Config{FailureModel: model})
		if err != nil {
			t.Errorf("model %q: unexpected error %v", model, err)
		}
		if got := draw(1.5, 100); math.IsNaN(got) || math.IsInf(got, 0) || got < 0 {
			t.Errorf("model %q: draw returned %v, want a finite non-negative value", model, got)
		}
	}

	// An unknown model must return an error and a zero-valued fallback draw.
	draw, err := ExtractFailureModel(configextractor.Config{FailureModel: "random"})
	if err == nil {
		t.Error("unknown model: expected an error, got nil")
	}
	if got := draw(1.5, 100); got != 0.0 {
		t.Errorf("unknown model: fallback draw returned %v, want 0.0", got)
	}
}

func TestDrawNumberWeibull(t *testing.T) {
	sessionLength := DrawNumberWeibull(1.5, 100)
	if math.IsNaN(sessionLength) || math.IsInf(sessionLength, 0) || sessionLength < 0 {
		t.Errorf("Weibull draw returned %v, want a finite non-negative value", sessionLength)
	}
}

func TestDrawNumberLognormal(t *testing.T) {
	sessionLength := DrawNumberLognormal(1.0, 0.5)
	if math.IsNaN(sessionLength) || math.IsInf(sessionLength, 0) || sessionLength < 0 {
		t.Errorf("lognormal draw returned %v, want a finite non-negative value", sessionLength)
	}
}

func TestExtractConnectivityFactor(t *testing.T) {
	cases := map[string]float64{
		"peer":       0.5,
		"peeper":     0.3,
		"benefactor": 0.7,
	}
	for profile, want := range cases {
		got, err := ExtractConnectivityFactor(configextractor.Config{NodeProfile: profile})
		if err != nil {
			t.Errorf("profile %q: unexpected error %v", profile, err)
		}
		if got != want {
			t.Errorf("profile %q: got %v, want %v", profile, got, want)
		}
	}

	got, err := ExtractConnectivityFactor(configextractor.Config{NodeProfile: "random"})
	if err == nil {
		t.Error("unknown profile: expected an error, got nil")
	}
	if got != 0.0 {
		t.Errorf("unknown profile: got %v, want 0.0", got)
	}
}

func TestComputeDowntimeFromSessionLength(t *testing.T) {
	// With a connectivity factor of 0.5 the downtime equals the session length.
	if got := computeDowntimeFromSessionLength(0.5, 10); got != 10 {
		t.Errorf("connectivity 0.5: got %v, want 10", got)
	}
	// A more available node (0.8) should be down for a quarter of its uptime.
	if got := computeDowntimeFromSessionLength(0.8, 10); math.Abs(got-2.5) > 1e-9 {
		t.Errorf("connectivity 0.8: got %v, want 2.5", got)
	}
}
