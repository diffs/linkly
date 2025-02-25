package server

import (
    "linkly/handlers"
    "regexp"

    "github.com/bwmarrin/discordgo"
    "github.com/sirupsen/logrus"
)

var linkDetection = regexp.MustCompile(`https?://\S+`)

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
    s.session.ShouldReconnectOnError = true
    s.session.Identify.Intents = discordgo.IntentsGuildMessages
    s.session.Identify.Intents |= discordgo.IntentsMessageContent

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

    link := linkDetection.FindString(m.Content)
    if link == "" {
        return
    }

    // Check if the message content matches the URL pattern
    response := s.handler.HandleLink(link)
    if response == nil {
        return
    }

    _, err := sess.ChannelMessageSendComplex(m.ChannelID, response)
    if err != nil {
        s.log.WithError(err).Error("failed to send message")
        return
    }
}
