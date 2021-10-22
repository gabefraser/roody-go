package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Token string

func init() {
	flag.StringVar(&Token, "t", "", "Bot token (required)")
	flag.Parse()

	if Token == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("Unable to create bot session.")
		os.Exit(1)
	}

	dg.AddHandler(onReady)
	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatal("Unable to open the connection to Discord.")
		os.Exit(1)
	}

	fmt.Println("Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func onReady(s *discordgo.Session, event * discordgo.Ready) {
	s.UpdateGameStatus(0, "i fucking dare you to 'invite' me to your server")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// A really hacky ugly way to fix nickname issues
	message := strings.Replace(m.Content, "!", "", -1)

	if strings.Contains(message, s.State.User.Mention()) && strings.Contains(message, "invite") {
		_, err := s.ChannelMessageSend(m.ChannelID, m.Author.Mention() + " let me yell at your friends with https://discord.com/api/oauth2/authorize?client_id=901108741318000692&permissions=0&scope=bot")
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if strings.Contains(message, s.State.User.Mention()) {
		_, err := s.ChannelMessageSend(m.ChannelID, m.Author.Mention() + " " + generateInsult())
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func generateInsult() string {
	insultClient := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://evilinsult.com/generate_insult.php", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Roody-Bot - https://github.com/gabefraser/roody-go")

	res, err := insultClient.Do(req)
	if err != nil {
		log.Fatal(nil)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
