package config


func DestinationRegistry(cfg Config) []string {

	// create a list of destination registries
	destinationRegistries := []string{}

	for _, repo := range cfg.Repos {
		for _, dst := range repo.Destinations {
			destinationRegistries = append(destinationRegistries, dst.Name)
		}
	}

	// remove duplicate values
	registryMap := make(map[string]bool)
	destinationRegistryList := []string{}

	for _, value := range destinationRegistries {
		if _, ok := registryMap[value]; !ok {
			registryMap[value] = true
			destinationRegistryList = append(destinationRegistryList, value)
		}
	}

	return destinationRegistryList
}
