package app

import "testing"

func TestSelfHostedProviderEndpoints(t *testing.T) {
	for _, raw := range []string{"http://127.0.0.1:2322", "http://[::1]:2322", "https://maps.example.test", "http://127.0.0.1:2322/tiles/{z}/{x}/{y}.png"} {
		if !validProviderEndpoint(raw) {
			t.Errorf("rejected configured provider %s", raw)
		}
	}
	for _, raw := range []string{"http://localhost:2322", "http://maps.example.test", "http://127.0.0.1.evil.test", "http://127.0.0.1@evil.test", "http://192.168.1.1", "file:///private", "https://user:pass@maps.example.test", "https://maps.example.test/#private"} {
		if validProviderEndpoint(raw) {
			t.Errorf("accepted insecure provider %s", raw)
		}
	}
}
