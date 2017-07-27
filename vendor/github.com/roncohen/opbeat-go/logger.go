package opbeat

// StdLogger is an interface that is implemented by the standard library
// `log.Logger`. Unfortunately it is no an interface in the standard library, so
// we are forced to include it here for interoperability.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}
