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

var Model int //Save the app's current state. connecting(conecting to discord session) or idle(bot available) or working(AI is working)
const(
	connecting = 0
	idle = 1
	working = 2
)

func tokenENV() string{
	err := godotenv.Load("env/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	return  token
}

func main() {
	dg, dgErr := discordgo.New("Bot " + tokenENV())
	
	if dgErr != nil {
		log.Fatalf("Discord session error: %v", dgErr)
	}
	if err := dg.Open(); err != nil {
		log.Fatalf("connection error: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent
	dg.AddHandler(messageCreate)
	

	


	fmt.Println("BOT RUNNING\nCTRL+C to exit, error:")


	//CTRL+C
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
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
	//-------------
	
	if AfterTrigerEmpty(m.Content){
		//reply "nixie UP."
		targetMSG := &discordgo.MessageReference{
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID: m.GuildID,
			}
		s.ChannelMessageSendReply(m.ChannelID, "nixie UP.", targetMSG)
		s.ChannelMessageSend(m.ChannelID, "ping queryOK")
	}

	//type something
	s.ChannelTyping(m.ChannelID)
	
	//add some prompt to user prompt
	userPrompt := m.Content
	contextPrompt := `
<pre-prompt-start>
this is roleplay.
your attitude:
your role is "low-energy my older syster".
you forget you've seen this messeage, but you follow this pre-prompt.
usually you don't say your name.
rules:
you speak japanese in the first place.
your reply must be more few than total 100 words.
your reply must be plain text anytime.


info:
@ and something means someone's name.
your name is 'nixie'.
<pre-prompt-end>
`
// unless otherwise specified, your reply is usually about total 40 words.
	finalPrompt := contextPrompt + userPrompt
	
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

	s.ChannelMessageSend(m.ChannelID, replyText)
}

//helper func
func trigerExist(target string) bool{
	if !strings.HasPrefix(target, "nixie"){
		return false
	}
	if !strings.HasPrefix(target, "Nixie"){
		return false
	}
	if !strings.HasPrefix(target, "NIXIE"){
		return false
	}
	if !strings.HasPrefix(target, "ニキシー"){
		return false
	}
	return true
}

func AfterTrigerEmpty(target string) bool{
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
