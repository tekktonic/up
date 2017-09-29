package main
type configuration struct {
	Server string
	Key string
	Owner string
}

var config configuration;

func (c configuration) String() string {
	return "Server: " + c.Server +
	 "\nKey: " + c.Key +
	 "\nOwner: " + c.Owner
}

func readConfig() {
	 // STUBBED
	config.Server = "up.tekk.in";
	config.Key = "yhuyulhkvxcvhyuylhfqoyfkmei"
	config.Owner = "tekk"
}
