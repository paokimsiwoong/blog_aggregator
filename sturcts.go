package main

import (
	"fmt"

	"github.com/paokimsiwoong/blog_aggregator/internal/config"
	"github.com/paokimsiwoong/blog_aggregator/internal/database"
)

// config.Config의 포인터를 저장하는 구조체
type state struct {
	ptrCfg *config.Config
	ptrDB  *database.Queries
}

// 명령어 한개의 정보를 저장하는 구조체
type command struct {
	name string
	args []string
}

// 사용가능한 명령어들을 저장한 구조체
type commands struct {
	// map: 명령어 이름 => 해당 명령어 handler 함수
	commandMap map[string]func(*state, command) error
}

// 새 명령어 handler를 commands 구조체에 저장하는 메소드
// func (c *commands) register(name string, f func(*state, command) error) error {
func (c *commands) register(name string, f func(*state, command) error) { // @@@ 과제 함수 시그니쳐 잘못 본 것 문제에 맞게 수정
	// if _, ok := c.commandMap[name]; ok {
	// 	return fmt.Errorf("command %s has already been registered: ", name)
	// }

	c.commandMap[name] = f

	// return nil
}

// commands 구조체에서 주어진 cmd를 찾아 실행하는 ㅁ세ㅗ드
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commandMap[cmd.name]
	if !ok {
		return fmt.Errorf("no such command : %s", cmd.name)
	}

	if err := f(s, cmd); err != nil {
		return fmt.Errorf("error running command %s : %w", cmd.name, err)
	}
	return nil
}
