package main

import(
	"fmt"
)

const cred = "123password"
const salt = "1f2606b6eec1654e"

func hardcoded1(password string) bool {
	if (password == "password") {
		return true;
	}
	return false;
}

func hardcoded2() {
	var pwd string = "password";
	var key string = "jb4HjJ8FV5j5d6XW";
	var hash string = "c4a15dd62973f33c54bcc002d8ce5d517901053b";
	var errorStr string = "Error: Not Found";
	apiCred := "key123"
	var notACred1, notACred2, notACred3 string;
	notACred1 = "Error: undefined reference";
	notACred2 = "Segmentation fault";
	notACred3 = "error";

	if hardcoded1(pwd) {
		fmt.Println(pwd);
	}
	if hardcoded1(key) {
		fmt.Println(2);
	}
	if hardcoded1(hash) {
		fmt.Println(2);
	}
	if hardcoded1("d99fce9480205c4b201fbc5fa80fd3232a4eefb6fba5cbefb5702d171ec14c33") {
		fmt.Println(3);
	}
	if hardcoded1(apiCred) {
		fmt.Println(apiCred);
	}

	fmt.Println(notACred1, notACred2, notACred3, errorStr);
}
