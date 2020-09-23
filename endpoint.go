package regip

// An Endpoint takes the database, first message of a new request chain, and the destination channel to output to, and returns
// An input channel to the endpoint, and a channel that signifies when the endpoint completes and can be GC'd
type Endpoint func(*DB, *string, *Message, chan *Message, chan *Message, *Logger) chan struct{}

// A SimpleEndpoint takes the database, first message, and an input and output channel (simpler to write)
type SimpleEndpoint func(*DB, *string, *Message, chan *Message, chan *Message, *Logger)

// Function to make it easy to compose Endpoints
func WrapEndpoint(e SimpleEndpoint) Endpoint {
	return func(database *DB, userid *string, first *Message, input, output chan *Message, lgr *Logger) chan struct{} {
		quit := make(chan struct{})
		go func() {
			e(database, userid, first, input, output, lgr) // Blocks until done / ready to quit
			quit <- struct{}{}
		}()
		return quit
	}
}
