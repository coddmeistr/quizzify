package config

type Service struct {
	Tests     Tests     `yaml:"tests" env-required:"true"`
	Questions Questions `yaml:"questions" env-required:"true"`
}

// TODO: write struct tags for all config fields

type Tests struct {
	ShortTextMaxLength int64 `yaml:"short_text_max_length" env-default:"60"`
	LongTextMaxLength  int64 `yaml:"long_text_max_length" env-default:"200"`
	ShortTextMinLength int64 `yaml:"short_text_min_length" env-default:"3"`
	LongTextMinLength  int64 `yaml:"long_text_min_length" env-default:"0"`
	MaxForCommonUser   int64 `yaml:"max_for_common_user" env-default:"10"`
	MaxForPremiumUser  int64
	MainImageByteSize  int64 `yaml:"main_image_byte_size" env-default:"4194304"`
}

type Questions struct {
	ShortTextMaxLength int64 `yaml:"short_text_max_length" env-default:"60"`
	LongTextMaxLength  int64 `yaml:"long_text_max_length" env-default:"200"`
	ShortTextMinLength int64 `yaml:"short_text_min_length" env-default:"3"`
	LongTextMinLength  int64 `yaml:"long_text_min_length" env-default:"0"`
	MaxForCommonUser   int64 `yaml:"max_for_common_user" env-default:"10"`
	MaxForPremiumUser  int64
}
