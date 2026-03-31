package Test

import "github.com/goravel/framework/contracts/ai"

type User struct {
}

func (r *User) Instructions() string {
	return ""
}

func (r *User) Messages() []ai.Message {
	return nil
}
