// Package diff helps build partial updates that only carry real changes.
package diff

// Changed returns requested when it is set and differs from current, and nil
// otherwise. A nil result means the field should be left out of the update.
func Changed[T comparable](current T, requested *T) *T {
	if requested == nil || *requested == current {
		return nil
	}
	return requested
}
