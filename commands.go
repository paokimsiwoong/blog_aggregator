package main

import (
	"errors"
	"fmt"
)

// login command 입력 시 실행되는 함수
func handlerLogin(s *state, cmd command) error {
	// login 뒤의 추가 명령어들은 cmd.args에 저장되어 있음
	// args의 길이가 1이 아니면 login <userName> 형태가 아니므로 에러
	if len(cmd.args) != 0 {
		return errors.New("the login handler expects a single argument, the username")
	}

	if err := s.ptrCfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error setting user : %w", err)
	}

	fmt.Println("the user has been set.")
	return nil
}
