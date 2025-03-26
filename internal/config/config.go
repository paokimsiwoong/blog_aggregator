// json 파일 읽기, 쓰기를 담당하는 패키지
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Config 구조체의 CurrentUserName string `json:"current_user_name"` 부분을 채우고 그 구조체를 json파일로 저장하는 함수
func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name

	if err := write(*c); err != nil {
		return fmt.Errorf("error writing file : %w", err)
	}

	return nil

	// @@@ 해답의 경우 return write(*c) 으로 바로 반환
}

// home 디렉토리에 있는 .gatorconfig.json 파일을 읽어와 Config 구조체에 저장하고 반환하는 함수
func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting file path : %w", err)
	}

	file, err := os.Open(filePath) // Open(name string) (*os.File, error)
	if err != nil {
		return Config{}, fmt.Errorf("error opening file : %w", err)
	}
	defer file.Close()

	var config Config

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("error decoding file : %w", err)
	}
	return config, nil
}

// config 파일 path 반환하는 함수
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home dir : %w", err)
	}

	return homeDir + "/" + configFileName, nil
	// @@@ 해답은 "path/filepath" 패키지를 사용해서 return filepath.Join(homeDir, configFileName), nil
}

// Config 구조체를 파일로 디스크에 저장하는 함수
func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshalling config : %w", err)
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting file path : %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing file (os.WriteFile) : %w", err)
	}
	// os.WriteFile 3번째 인자는 permission
	// ex: -rw-r--r--는 0644
	// 0(owner)(group)(guest), x는 1 w는 2 r는 4 ==> rwx = 4 + 2 + 1, rw- = 4 + 2, r-x = 4 + 1, ...

	// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
	// @@@ 해답 비교 (json.NewEncoder(file) 활용해서 Read와 write 함수가 구조가 거의 동일)
	// fullPath, err := getConfigFilePath()
	// if err != nil {
	// 	return err
	// }

	// file, err := os.Create(fullPath)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// encoder := json.NewEncoder(file)
	// err = encoder.Encode(cfg)
	// if err != nil {
	// 	return err
	// }
	// return nil
	// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

	return nil
}
