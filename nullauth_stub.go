//go:build !nullauth

package jasmine

func NullAuth(provider AuthProvider) AuthProvider {
	return provider
}
