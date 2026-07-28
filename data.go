case "stop":
	fs := flag.NewFlagSet("stop", flag.ContinueOnError)
	fs.IntVar(&f.Port, "port", 8080, "Port used by the Replay Engine server")
	fs.StringVar(
		&f.Config,
		"config",
		"",
		"Path to a RunConfig JSON file",
	)
	fs.SetOutput(os.Stderr)
	_ = fs.Parse(remaining)