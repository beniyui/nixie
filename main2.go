package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/bwmarrin/discordgo"
)

func tokenENV() string{
	err := godotenv.Load("env/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	return  token
}

const (
	OllamaModel  = "gemma4:e2b"
	OllamaURL    = "http://localhost:11434/api/generate"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Think bool   `json:"think"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func main() {
	dg, err := discordgo.New("Bot " + tokenENV())
	if err != nil {
		log.Fatalf("Discord session error: %v", err)
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent

	err = dg.Open()
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}

	fmt.Println("BOT RUNNING\nCTRL+C to exit, error:")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}



func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ingnore herself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// bot triger
	prefix1 := "nixie"
	if len(m.Content) <= len(prefix1) || m.Content[:len(prefix1)] != prefix1 {
		return
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
	replyText, err := queryOllama(finalPrompt)
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



func queryOllama(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  OllamaModel,
		Prompt: prompt,
		Stream: false,
		Think: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Response, nil
}
