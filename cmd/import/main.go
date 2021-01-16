package main

func main() {
	arg := parseArgs()

	if arg.authorize {
		authorize(arg)
		return
	}

	importData(arg)
}
