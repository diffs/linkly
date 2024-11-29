package server

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"linkly/handlers"
	"regexp"
)

var linkDetection = regexp.MustCompile(`https?://[^\s]+`)

type Server interface {
	Start() error
}

type server struct {
	session *discordgo.Session
	log     *logrus.Logger
	handler handlers.Handler
}

func NewServer(session *discordgo.Session, handler handlers.Handler) Server {
	return &server{
		session: session,
		log:     logrus.New(),
		handler: handler,
	}
}

func (s *server) Start() error {
	s.session.AddHandler(s.messageCreate)

	s.session.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := s.session.Open()
	if err != nil {
		s.log.WithError(err).Error("failed to open ws connection to Discord")
		return err
	}

	if s.session.State != nil {
		s.log.WithField("user", s.session.State.User.ID).Info("Connected to discord!")
	}
	return nil
}

func (s *server) messageCreate(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == sess.State.User.ID {
		return
	}

	fmt.Println(m.Content)
	// Check if the message content matches the URL pattern
	if linkDetection.MatchString(m.Content) {
		if response := s.handler.HandleLink(m.Content); response != nil {
			if _, err := sess.ChannelMessageSendComplex(m.ChannelID, response); err != nil {
				s.log.WithError(err).Error("failed to send message")
				return
			}
		}
	}
}
