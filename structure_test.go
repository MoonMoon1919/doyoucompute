package doyoucompute

import "testing"

func testOperation[T Structurer, R any](
	t *testing.T,
	structUnderTest T,
	operation func(T) (R, error),
	errorMessage string,
	comparisonFunc func(T, R, *testing.T),
) {
	res, err := operation(structUnderTest)

	checkErrors(errorMessage, err, t)
	if errorMessage != "" {
		return // bail out before validation for tests w/ non-null errors
	}

	comparisonFunc(structUnderTest, res, t)
}
