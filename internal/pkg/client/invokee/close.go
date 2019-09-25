package invokee

// Close closes connection.
func Close() (e error) {
	e = conn.Close()

	return
}
