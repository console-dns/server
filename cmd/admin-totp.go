package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/console-dns/server/pkg/content/settings"
	"github.com/console-dns/server/pkg/utils"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
)

var totpCmd = &cobra.Command{
	Use:   "totp",
	Short: "TOTP 验证码管理",
}

var totpResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "重置当前 TOTP 验证码",
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
		fmt.Printf(`
###########################################################
#                     WARNING                             #
#   Turning off OTP will result in reduced security.      #
#                                                         #
###########################################################

`)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press Y to continue: ")
		text, _, _ := reader.ReadLine()
		if string(text) != "Y" && string(text) != "y" {
			fmt.Printf("The wrong key (%s) was entered，Exiting...\n", text)
			os.Exit(1)
		}
		staticCfg.Auth.TotpSecret = ""
		err = utils.AutoMarshal(configPath, staticCfg)
		if err != nil {
			fmt.Printf("Failed to write config: %v\n", err)
			os.Exit(1)
		}
	},
}

var totpCreateCmd = &cobra.Command{
	Use:   "generate",
	Short: "创建新的 TOTP 验证码",
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
		if staticCfg.Auth.TotpSecret != "" {
			fmt.Println("Totp secret already set.")
			os.Exit(1)
		}
		generate, _ := totp.Generate(totp.GenerateOpts{
			Issuer:      "DNS Console 2FA",
			AccountName: staticCfg.Auth.Username,
		})
		fmt.Printf(`
###########################################################
#                     WARNING                             #
#    You need to stop the current instance first.         #
#                                                         #
#  Please use your password management tool e.g.          #
#  Google Authenticator to scan the QR code above to add  #
#  an account, or you can use the link below to add       #
#  a TOTP account.                                        #
###########################################################
`)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press Y to continue: ")
		text, _, _ := reader.ReadLine()
		if string(text) != "Y" && string(text) != "y" {
			fmt.Printf("The wrong key (%s) was entered，Exiting...\n", text)
			os.Exit(1)
		}
		qrterminal.GenerateWithConfig(generate.String(), qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout,
			BlackChar: qrterminal.BLACK,
			WhiteChar: qrterminal.WHITE,
			QuietZone: 1,
		})
		fmt.Printf(`
totp secret: %s
totp url   : %s
`, generate.Secret(), generate.URL())
		ok := false
		for i := 1; i <= 3; i++ {
			fmt.Print("Enter your OTP code: ")
			totpCode, _, _ := reader.ReadLine()
			if string(totpCode) == "" {
				os.Exit(1)
			}
			if totp.Validate(string(totpCode), generate.Secret()) {
				ok = true
				break
			}
			fmt.Printf("verification code is incorrect [%d/3].\n", i)
		}
		if !ok {
			os.Exit(1)
		}
		staticCfg.Auth.TotpSecret = generate.Secret()
		err = utils.AutoMarshal(configPath, staticCfg)
		if err != nil {
			fmt.Printf("Failed to write config: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	totpCmd.AddCommand(totpResetCmd)
	totpCmd.AddCommand(totpCreateCmd)
	adminCmd.AddCommand(totpCmd)
}
