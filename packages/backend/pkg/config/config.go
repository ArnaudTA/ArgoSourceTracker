package config

var Global Config

func LoadGlobal() error {
	return Load(&Global)
}
