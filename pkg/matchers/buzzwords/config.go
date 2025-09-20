package buzzwords

import (
	"fmt"
	"regexp"
	"strings"

	matcher "github.com/br0-space/bot-matcher"
)

type Buzzword struct {
	Trigger string `mapstructure:"trigger"`
	Pattern string `mapstructure:"pattern"`
	Reply   string `mapstructure:"reply"`
}

func (b Buzzword) Matches(text string) bool {
	if b.Pattern != "" {
		pattern := fmt.Sprintf(`(?i)^%s$`, b.Pattern)

		return regexp.MustCompile(pattern).MatchString(text)
	}

	return strings.EqualFold(b.Trigger, text)
}

type Config struct {
	matcher.Config

	Buzzwords []Buzzword `mapstructure:"buzzwords"`
}

func (c Config) GetPattern() string {
	var patterns []string

	for _, buzzword := range c.Buzzwords {
		if buzzword.Pattern != "" {
			patterns = append(patterns, buzzword.Pattern)
		} else {
			patterns = append(patterns, buzzword.Trigger)
		}
	}

	return strings.Join(patterns, "|")
}

func (c Config) GetTrigger(match string) string {
	match = strings.TrimSpace(match)

	for _, buzzword := range c.Buzzwords {
		if buzzword.Matches(match) {
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
