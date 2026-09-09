package posts

type integer interface {
	~int | ~int64
}

func resolveMax[T integer](requested, defaultValue, maximum T) T {
	if requested <= 0 {
		requested = defaultValue
	}

	return min(requested, maximum)
}

func resolveMinMax[T integer](requested, defaultValue, minimum, maximum T) T {
	if requested <= 0 {
		requested = defaultValue
	}

	return min(max(requested, minimum), maximum)
}

// apiResults widens the request to the endpoint minimum: the caller may want 3
// posts, but the API rejects anything below apiMinResults, so the surplus is
// trimmed after the response arrives.
func apiResults(requested int) int {
	return max(apiMinResults, min(requested, 100))
}
