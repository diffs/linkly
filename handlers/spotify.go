package handlers

import (
	"context"
	"fmt"
	"linkly/handlers/utils"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const spotifyTokenExpiry = time.Hour

type SpotifyHandler struct {
	lastCall      time.Time
	client        *spotify.Client
	ctx           context.Context
	youtubeAPIKey string
}

func NewSpotifyHandler(ctx context.Context, client *spotify.Client, youtubeAPIKey string) Handler {
	return &SpotifyHandler{client: client, ctx: ctx, youtubeAPIKey: youtubeAPIKey, lastCall: time.Now()}
}

func (s *SpotifyHandler) HandleLink(link string) *discordgo.MessageSend {
	if !strings.Contains(link, "/track/") {
		return nil
	}
	splat := strings.Split(link, "/")

	trackID := strings.Split(splat[len(splat)-1], "?")[0]
	if s.isSpotifyClientTokenExpired() {
		fmt.Println("renewing token...")
		if err := s.renewSpotifyClient(); err != nil {
			fmt.Println(err.Error())
			return nil
		}
	}
	track, err := s.client.GetTrack(s.ctx, spotify.ID(trackID))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// Create the embed
	embed := &discordgo.MessageEmbed{
		Title:       track.Name,
		Description: "by " + track.Artists[0].Name,
		URL:         strings.Split(link, "?")[0],
		Color:       0x1DB954, // Spotify green
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Album",
				Value:  track.Album.Name,
				Inline: true,
			},
			{
				Name:   "Release Date",
				Value:  track.Album.ReleaseDate,
				Inline: true,
			},
		},
	}

	// Let's check for other links, in case someone doesn't wanna pay for Spotify...
	ytLink, _ := utils.QueryYoutube(track.ExternalIDs["isrc"], s.youtubeAPIKey)

	// We are going under the assumption that if there is a YouTube link returned from querying the ISRC, the first one is the correct link to the song
	if ytLink != nil && len(ytLink.Items) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Other Links",
			Value:  "[YouTube](https://youtu.be/" + ytLink.Items[0].ID.VideoID + ")",
			Inline: true,
		})
	}

	// If there are images, set the thumbnail
	if len(track.Album.Images) > 0 {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: track.Album.Images[0].URL,
		}
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}
}

func (s *SpotifyHandler) isSpotifyClientTokenExpired() bool {
	return time.Now().After(s.lastCall.Add(spotifyTokenExpiry))
}

func (s *SpotifyHandler) renewSpotifyClient() error {
	tk, err := s.client.Token()
	if err != nil {
		return err
	}

	httpClient := spotifyauth.New().Client(s.ctx, tk)
	s.client = spotify.New(httpClient)
	s.lastCall = time.Now()
	return nil
}
