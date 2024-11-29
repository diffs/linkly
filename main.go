package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ianschenck/envflag"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"linkly/handlers"
	"linkly/server"
	"os"
	"os/signal"
	"syscall"
)

var (
	authToken           = envflag.String("DISCORD_AUTH_TOKEN", "", "Discord auth token")
	spotifyClientID     = envflag.String("SPOTIFY_ID", "", "Spotify api client ID")
	spotifyClientSecret = envflag.String("SPOTIFY_SECRET", "", "Spotify api client secret")
)

func main() {
	envflag.Parse()

	discord, err := discordgo.New("Bot " + *authToken)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     *spotifyClientID,
		ClientSecret: *spotifyClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		panic(err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	spotifyClient := spotify.New(httpClient)
	spotifyHandler := handlers.NewSpotifyHandler(ctx, spotifyClient)
	hm := handlers.HandleManager{Handlers: map[string]handlers.Handler{"open.spotify.com": spotifyHandler}}

	ss := server.NewServer(discord, &hm)
	if err := ss.Start(); err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	_ = discord.Close()
}
