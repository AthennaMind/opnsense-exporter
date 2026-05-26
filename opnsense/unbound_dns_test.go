package opnsense

import "testing"

func TestParseUnboundOverview(t *testing.T) {
	const endpoint EndpointPath = "/api/unbound/diagnostics/stats"

	t.Run("populated DNSSEC counters", func(t *testing.T) {
		var resp unboundDNSStatusResponse
		resp.Data.Time.Up = "12345.678"
		resp.Data.Num.Answer.Bogus = "3"
		resp.Data.Num.Answer.Secure = "42"

		got, err := parseUnboundOverview(&resp, endpoint)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.UptimeSeconds != 12345.678 {
			t.Errorf("UptimeSeconds = %v, want 12345.678", got.UptimeSeconds)
		}
		if got.AnswerBogusTotal != 3 {
			t.Errorf("AnswerBogusTotal = %d, want 3", got.AnswerBogusTotal)
		}
		if got.AnswerSecureTotal != 42 {
			t.Errorf("AnswerSecureTotal = %d, want 42", got.AnswerSecureTotal)
		}
		if got.QueryTypes == nil || got.AnswerRcodes == nil {
			t.Error("QueryTypes / AnswerRcodes maps should be initialized")
		}
	})

	t.Run("zero DNSSEC counters", func(t *testing.T) {
		var resp unboundDNSStatusResponse
		resp.Data.Time.Up = "0"
		resp.Data.Num.Answer.Bogus = "0"
		resp.Data.Num.Answer.Secure = "0"

		got, err := parseUnboundOverview(&resp, endpoint)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.AnswerBogusTotal != 0 || got.AnswerSecureTotal != 0 {
			t.Errorf("expected zeroed counters, got bogus=%d secure=%d",
				got.AnswerBogusTotal, got.AnswerSecureTotal)
		}
	})

	t.Run("empty uptime returns error", func(t *testing.T) {
		var resp unboundDNSStatusResponse
		// Extended Statistics disabled: Time.Up empty, Num.Answer fields zero-value.
		// Existing behavior should fail loudly so operators see they need to enable it.
		_, err := parseUnboundOverview(&resp, endpoint)
		if err == nil {
			t.Fatal("expected error for empty uptime, got nil")
		}
		if err.Endpoint != string(endpoint) {
			t.Errorf("error endpoint = %q, want %q", err.Endpoint, string(endpoint))
		}
	})

	t.Run("empty bogus surfaces parse error", func(t *testing.T) {
		var resp unboundDNSStatusResponse
		resp.Data.Time.Up = "1.0"
		// Bogus left empty — must surface a parse error so the operator sees
		// the upstream API didn't populate the DNSSEC fields.
		_, err := parseUnboundOverview(&resp, endpoint)
		if err == nil {
			t.Fatal("expected error for empty bogus counter, got nil")
		}
	})
}
