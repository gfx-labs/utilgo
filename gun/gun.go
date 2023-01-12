package gun

import (
	"fmt"
	"os"
	"path"

	"gfx.cafe/util/go/gun/gunyaml"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
)

func Load(i any) {
	LoadPrefix(i, "")
}

func LoadPrefix(i any, prefix string) {
	yamlDecoder := gunyaml.New()
	dotenvDecoder := aconfigdotenv.New()
	fileName := "config"
	if prefix != "" {
		fileName = prefix
	}
	homeDir, _ := os.UserHomeDir()

	loader := aconfig.LoaderFor(i, aconfig.Config{
		AllowUnknownFields: true,
		AllowUnknownEnvs:   true,
		AllowUnknownFlags:  true,
		SkipFlags:          true,
		DontGenerateTags:   true,
		MergeFiles:         true,
		EnvPrefix:          prefix,
		FlagPrefix:         prefix,
		Files: []string{
			fmt.Sprintf("/%s.yml", fileName),
			fmt.Sprintf("/%s.yaml", fileName),
			fmt.Sprintf("/%s.json", fileName),
			fmt.Sprintf("/config/%s.yml", fileName),
			fmt.Sprintf("/config/%s.yaml", fileName),
			fmt.Sprintf("/config/%s.json", fileName),
			path.Join(homeDir, fmt.Sprintf(".gfx/%s.yml", fileName)),
			path.Join(homeDir, fmt.Sprintf(".gfx/%s.yaml", fileName)),
			path.Join(homeDir, fmt.Sprintf(".gfx/%s.json", fileName)),
			".env",
		},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": yamlDecoder,
			".yml":  yamlDecoder,
			".json": yamlDecoder,
			".env":  dotenvDecoder,
		},
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
}
