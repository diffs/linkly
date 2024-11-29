package handlers

import (
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type HandleManager struct {
	Handlers map[string]Handler
}

type Handler interface {
	HandleLink(link string) *discordgo.MessageSend
}

func (hm *HandleManager) HandleLink(link string) *discordgo.MessageSend {
	// Get the domain of the link
	domain := strings.ToLower(getDomain(link))
	handler, ok := hm.Handlers[domain]
	if !ok {
		return nil
	}

	return handler.HandleLink(link)
}

// getDomain extracts the domain from a given link URL string.
func getDomain(link string) string {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}
