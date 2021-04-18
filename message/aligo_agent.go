package message

// aligoAgent is struct that agent API about aligo including sending message, get message list, etc ...
type aligoAgent struct {
	apiKey, id, sender string
}

func AligoAgent(apiKey, id, sender string) *aligoAgent {
	return &aligoAgent{
		apiKey: apiKey,
		id:     id,
		sender: sender,
	}
}
