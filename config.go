package main
type configuration struct {
	Server string
	Key string
	Owner string
}

var config configuration;

func readConfig() {
	 // STUBBED
	config.Server = "up.tekk.in";
	config.Key = "yhuyulhkvxcvhyuylhfqoyfkmei"
	config.Owner = "tekk"
}
