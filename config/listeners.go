package config

type Listeners struct {
	Listeners []Listener `listeners`
}

type Listener struct {
	Protocol string
	Port     int
}
