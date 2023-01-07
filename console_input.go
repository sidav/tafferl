package main

func readKey() string {
	return cw.ReadKey()
}

func readKeyAsync() string {
	return cw.ReadKeyAsync(10)
}
