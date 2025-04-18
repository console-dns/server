package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/alexedwards/argon2id"
	"github.com/console-dns/server/pkg/content/settings"
	"github.com/console-dns/server/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var passwordCmd = &cobra.Command{
	Use:   "pass",
	Short: "设置管理员密码",
	Run: func(cmd *cobra.Command, args []string) {
		if configPath == "" {
			fmt.Println("未指定配置文件")
			os.Exit(1)
		}
		staticCfg, err := settings.FromStaticConfig(configPath)
		if err != nil {
			fmt.Printf("Failed to read config: %v\n", err)
			os.Exit(1)
		}
		var password string
		for i := 1; i <= 3; i++ {
			password, err = credentials()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			password, err = argon2id.CreateHash(password, argon2id.DefaultParams)
			if err == nil {
				break
			} else {
				fmt.Println(err.Error())
			}

		}
		if err != nil {
			fmt.Printf("The password is not updated")
			os.Exit(1)
		}
		staticCfg.Auth.Password = "argon2:" + password
		err = utils.AutoMarshal(configPath, staticCfg)
		if err != nil {
			fmt.Printf("Failed to write config: %v\n", err)
			os.Exit(1)
		}
	},
}

func credentials() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	println()
	fmt.Print("Enter Password Again: ")
	bytePasswordAgain, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	println()

	password := string(bytePassword)
	passwordAgain := string(bytePasswordAgain)
	if password != passwordAgain {
		return "", errors.New("Passwords don't match")
	}
	if len(password) < 8 {
		return "", errors.New("Password must be at least 8 characters")
	}
	return strings.TrimSpace(password), nil
}

func init() {
	adminCmd.AddCommand(passwordCmd)
}
