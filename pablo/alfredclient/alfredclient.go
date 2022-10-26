package alfredclient

const (
	authURL = "/auth"
)

type Implementation struct { // not a huge fan of this name, but as far as this package goes… it's the implementation of the main thing
	addr string
	auth AuthResponse
}

func New(addr string) *Implementation {
	return &Implementation{addr: addr}
}
