package main

import (
	"math/rand"
	"strings"

	"github.com/mattn/go-xmpp"
	"github.com/spf13/viper"
)

func (bot *botClient) answer(chat xmpp.Chat) {

	typeXMPP := getTypeChat(chat.Remote)
	remote := chat.Remote
	to := strings.Split(chat.Remote, "@")[0]
	if typeXMPP == "groupchat" {
		if strings.Contains(chat.Remote, "/") {
			t := strings.Split(chat.Remote, "/")
			remote = t[0]
			to = t[1]
		}
	}

	bot.chats <- xmpp.Chat{
		Remote: remote,
		Type:   typeXMPP,
		Text:   to + ": " + bot.prepareAnswer(chat.Text, to, chat.Remote),
	}
	bot.nbXMPPAnswers++
}

func (bot *botClient) prepareAnswer(text, short, remote string) string {
	question := strings.TrimSpace(text[5:]) // remove '/cds ' or 'cds, '

	switch question {
	case "help":
		return help()
	case "cds2xmpp status":
		if bot.isAdmin(remote) {
			return bot.getStatus()
		}
		return "forbidden for you " + remote
	case "ping":
		return "pong"
	default:
		return random()
	}

}

func help() string {
	out := `
Begin conversation with "cds," or "/cds"

Simple request: "cds, ping"

/cds cds2xmpp status (for admin only)

`

	return out + viper.GetString("more_help")
}

func (bot *botClient) isAdmin(r string) bool {
	for _, a := range bot.admins {
		if strings.HasPrefix(r, a) {
			return true
		}
	}
	return false
}

func random() string {
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
		"Nooooo",
	}
	return answers[rand.Intn(len(answers))]
}
