package main;

import ("fmt")
func auth(secret string) string {

	fmt.Println("Trying to authenticate " + secret + " against " + config.Key)
	if (secret == config.Key) {
		return "";
	}

	return "Not authorized"
}
