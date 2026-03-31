package agents

import "github.com/goravel/framework/contracts/ai"

type Test struct {
}

func (r *Test) Instructions() string {
	return "You are a test agent."
}

func (r *Test) Messages() []ai.Message {
	return []ai.Message{}
}
