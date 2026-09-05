package config

import (
	"strings"
	"testing"
)

func TestValidateForAPI(t *testing.T) {
	cases := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{"empty is rejected", "", true},
		{"whitespace-only is rejected", "        ", true},
		{"too short is rejected", "short", true},
		{"exactly min length is accepted", strings.Repeat("x", minSessionSecretLen), false},
		{"long secret is accepted", "a-perfectly-reasonable-session-secret", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Config{SessionSecret: tc.secret}
			err := c.ValidateForAPI()
			if tc.wantErr && err == nil {
				t.Errorf("expected error for secret %q, got nil", tc.secret)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error for secret %q, got %v", tc.secret, err)
			}
		})
	}
}
