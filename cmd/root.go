package cmd

import (
	"custodyEthereum/configs"
	"custodyEthereum/internal/jwtHelper"
	"custodyEthereum/pkg/encryptedStore"
	"custodyEthereum/pkg/server"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "raccoon-backend",
	Short: "Run the BackEnd for Raccoon Fantasy",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		startCustodyServer()
	},
}

var (
	cfgFile string
	sslFlag bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config_path", "./configs/local.yaml", "config file (default is $HOME/v1/configs/local.yaml)")
	rootCmd.PersistentFlags().BoolVar(&sslFlag, "ssl_flag", false, "Disables SSL (default is false)")

}

func initConfig() {

	vp := viper.New()

	if cfgFile != "" {
		cfgFileList := strings.Split(cfgFile, "/")
		path := strings.Join(cfgFileList[:len(cfgFileList)-1], "/")
		file := strings.Split(cfgFileList[len(cfgFileList)-1], ".")

		vp.SetConfigName(file[0])
		vp.SetConfigType(file[1])
		vp.AddConfigPath(path)
	}

	vp.BindPFlag("server.ssl_dis", rootCmd.PersistentFlags().Lookup("ssl_flag"))

	configs.StartGlobalConfig(vp)
}

func startCustodyServer() {
	log.Println("Server settings:")
	log.Println(configs.GlobalViper.GetString("server.host"))
	log.Println(configs.GlobalViper.GetString("server.port"))

	firstStart := true

	s := server.Server{
		AvailableStores: make(map[string]bool),
	}

	files, err := os.ReadDir(configs.GlobalViper.GetString("server.basepath") + "/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		s.AvailableStores[file.Name()] = true
		if file.Name() == configs.GlobalViper.GetString("server.default-storage") {
			firstStart = false
		}
	}

	if firstStart {
		log.Println("First start detected, creating default storage, and account")

		//Create root account
		rootToken, err := jwtHelper.GenerateAccessToken(false, "root")
		if err != nil {
			log.Fatal(err)
		}

		//Create default storage
		shares := encryptedStore.CreateNewBoxStore(configs.GlobalViper.GetString("server.default-storage"), 3, 5, "")

		var textShares []string
		for _, share := range shares {
			sh, err := share.Open()
			if err != nil {
				log.Println(err)
			}
			textShares = append(textShares, base64.StdEncoding.EncodeToString(sh.Bytes()))
		}

		//Show credentials
		fmt.Printf(strings.Repeat("=", 12) + "\n")
		fmt.Printf("Store Name: %s\n", configs.GlobalViper.GetString("server.default-storage"))
		fmt.Printf("Store Status: %s\n", "Locked")
		fmt.Printf(strings.Repeat("*", 12) + "\n")
		fmt.Printf("Root Token: %s\n", rootToken)
		fmt.Printf(strings.Repeat("*", 12) + "\n")
		fmt.Printf("Store Keys:\n")
		fmt.Printf("%-12s %-20s\n", "Share", "Key")
		for i, share := range textShares {
			fmt.Printf("%-12s %-20s\n", "Share "+strconv.Itoa(i), share)
		}
		fmt.Printf(strings.Repeat("=", 12) + "\n")

		rootToken = ""
		textShares = nil
	}

	r := gin.Default()

	r.GET("initialize", jwtHelper.JwtAuthGetRoleMiddleware(), s.NewStore())

	r.POST("unlock", jwtHelper.JwtAuthGetRoleMiddleware(), s.Unlock(false))
	r.POST("reloadStore", jwtHelper.JwtAuthGetRoleMiddleware(), s.Unlock(true))
	r.POST("addSecret", jwtHelper.JwtAuthGetRoleMiddleware(), s.AddSecret())
	r.POST("remSecret", jwtHelper.JwtAuthGetRoleMiddleware(), s.RemoveSecret())
	r.POST("updSecret", jwtHelper.JwtAuthGetRoleMiddleware(), s.UpdateSecret())
	//r.POST("signTransaction", jwtHelper.JwtAuthGetRoleMiddleware(), s.SignTransaction)

	err = r.Run(":8080")
}
