package globals

import (
	"github.com/pelletier/go-toml"
)

var (
	BotSettings *toml.Tree
	AccessToken string
	FullBase    [][]string
	HasAnswers  bool = false
)
