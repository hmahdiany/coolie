package config


func DestinationRegistry(cfg Config) map[string]string {

	// create a list of destination registries
	destinationRegistries := map[string]string{}

	for _, repo := range cfg.Repos {
		for _, dst := range repo.Destinations {
			destinationRegistries[dst.Name] = dst.Address
		}
	}

	return destinationRegistries
}
