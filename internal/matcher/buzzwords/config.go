package buzzwords

import (
	"fmt"
	"github.com/br0-space/bot/internal/matcher/abstract"
	"regexp"
	"strings"
)

type Config struct {
	abstract.Config
	Buzzwords []Buzzword `mapstructure:"buzzwords"`
}

type Buzzword struct {
	Trigger string `mapstructure:"trigger"`
	Pattern string `mapstructure:"pattern"`
	Reply   string `mapstructure:"reply"`
}

func (c Config) GetPattern() string {
	var patterns []string
	for _, buzzword := range c.Buzzwords {
		if buzzword.Pattern != "" {
			patterns = append(patterns, buzzword.Pattern)
		} else {
			patterns = append(patterns, string(buzzword.Trigger))
		}
	}

	return strings.Join(patterns, "|")
}

func (c Config) GetTrigger(match string) string {
	for _, buzzword := range c.Buzzwords {
		if buzzword.Pattern != "" && regexp.MustCompile(buzzword.Pattern).MatchString(match) {
			return buzzword.Trigger
		}
		if buzzword.Trigger == match {
			return buzzword.Trigger
		}
	}

	return ""
}

func (c Config) GetReply(match string) (string, error) {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?i)(\W)%s(\W)`, match))
	var reply string

	for _, buzzword := range c.Buzzwords {
		if buzzword.Pattern != "" && pattern.MatchString(buzzword.Pattern) {
			reply = buzzword.Reply
			break
		}
		if buzzword.Trigger == match {
			reply = buzzword.Reply
			break
		}
	}

	if reply == "" {
		return "", fmt.Errorf("no reply found for match %s", match)
	}

	return reply, nil
}
