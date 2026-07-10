if len(os.Args) > 1 {
	switch os.Args[1] {
	case "version", "--version", "-v":
		fmt.Printf("rre version %s\n", version)
		return
	}
}