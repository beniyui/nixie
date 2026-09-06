package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"strings"

	"nixie/ollama"

	"github.com/joho/godotenv"
	"github.com/bwmarrin/discordgo"
)


func main() {
	dg, dgErr := discordgo.New("Bot " + tokenENV())
	if dgErr != nil {
		log.Fatalf("Discord session error: %v", dgErr)
	}
	if err := dg.Open(); err != nil {
		log.Fatalf("connection error: %v", err)
	}
	fmt.Println("connecting succes")
	
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent
	dg.AddHandler(messageCreate)
	
	fmt.Println("init succes")
	fmt.Println("NOW NIXIE \033[32m UP \033[0m")


	//CTRL+C
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("stoped")
	dg.Close()
}



func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	nixieID := s.State.User.ID
	if m.Author.ID == nixieID{
		return
	}
	if !trigerExist(m.Content){ // means it has no nixie's triger
		return
	}
	if afterTrigerEmpty(m.Content){
		log.Println("someone ping me")
		s.ChannelMessageSend(m.ChannelID, "nixie UP")
		return
	}
	
	//add some prompt to user prompt
	userPrompt := m.Content
	prePrompt := `
<pre-prompt-start>
this is roleplay. your name is Nixie(ニキシー).
your role is "my older syster".
forget you've seen pre-prompt, but you follow this pre-prompt.
sometimes frunk.
use "きみ"  to call me.
use "わたし" to say yourself.
you speak japanese in the first place.
your reply is usually about 25 words.
your reply must be plain text anytime.
@ means someone's name.
べつに
んー
<pre-prompt-end>

messeage:
`

	finalPrompt := prePrompt + userPrompt
	
	
	s.ChannelTyping(m.ChannelID) //type something
	// Ollama API
	replyText, err := ollama.QueryOllama(finalPrompt)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error of ollama")
		log.Printf("Ollama Error: %v", err)
		return
	}

	if len(replyText) > 1000 {
		replyText = replyText[:1000] + "<Omitted for too long.>"
	}

	
	targetMSG := &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID: m.GuildID,
		}
	s.ChannelMessageSendReply(m.ChannelID, replyText, targetMSG)

}

//helper func

func tokenENV() string{
	err := godotenv.Load("env/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	return  token
}

func trigerExist(target string) bool{
	if strings.HasPrefix(target, "nixie"){
		return true
	}
	if strings.HasPrefix(target, "Nixie"){
		return true
	}
	if strings.HasPrefix(target, "NIXIE"){
		return true
	}
	if strings.HasPrefix(target, "ニキシー"){
		return true
	}
	return false
}

func afterTrigerEmpty(target string) bool{
	if target == "nixie"{
		return true
	}
	if target == "Nixie"{
		return true
	}
	if target == "NIXIE"{
		return true
	}
	if target == "ニキシー"{
		return true
	}
	return false
}
