package freerdp

func Try(f func()) {
	defer func() { recover() }()
	f()
}
