package assigns

// Report passes a result of the assignment to the assignee.
func Report(invkID string, rst interface{}) (ok bool) {
	var (
		a assign
	)

	if a, ok = assignments[invkID]; !ok {
		return
	}

	a.holder <- rst

	delete(assignments, invkID)

	return
}
