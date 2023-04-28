package cmd

import (
	"custodyEthereum/configs"
	"custodyEthereum/internal/jwtHelper"
	"custodyEthereum/pkg/encryptedStore"
	"custodyEthereum/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config_path", "./v1/configs/local.yaml", "config file (default is $HOME/v1/configs/local.yaml)")
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
		shares := encryptedStore.CreateNewStore(configs.GlobalViper.GetString("server.default-storage"), 3, 5)

		//Show credentials
		log.Println("Root token: " + rootToken)
		log.Println("Default storage: " + configs.GlobalViper.GetString("server.default-storage"))
		log.Println("Shares: ", shares)
	}

	r := gin.Default()

	r.GET("initialize", jwtHelper.JwtAuthMiddleware([]string{"root"}), s.Initialize())

	r.POST("unlock", jwtHelper.JwtAuthMiddleware([]string{"root"}), s.Unlock())
	//r.POST("sign", jwtHelper.JwtAuthMiddleware([]string{"root"}), s.Sign)
	err = r.Run(":8080")
}
