package configs

import "github.com/spf13/viper"

var GlobalViper *viper.Viper

func StartGlobalConfig(vp *viper.Viper) {
	GlobalViper = initConfig(vp)
}

func initConfig(vp *viper.Viper) *viper.Viper {

	configKey := "raccoon/raccoon-be.yaml"

	/*consulConfig := api.DefaultConfig()
	consulConfig.Address = "consul.cluster.internal:80"
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		panic(err)
	}*/

	// Set default Server
	vp.SetDefault("server.host", "0.0.0.0")
	vp.SetDefault("server.port.https", 8080)
	vp.SetDefault("server.port.http", 7070)
	vp.SetDefault("server.readTimeout", 10)
	vp.SetDefault("server.writeTimeout", 10)
	vp.SetDefault("server.ssl_dis", false)
	// Set default jwt
	vp.SetDefault("jwt.accessprivkeypath", "access-private.pem.pem")
	vp.SetDefault("jwt.accesspubkeypath", "access-public.pem")

	// Set SSL default
	vp.SetDefault("ssl.keypath", "domain.key")
	vp.SetDefault("ssl.certpath", "domain.crt")

	// Set default Kafka
	vp.SetDefault("kafka.host", "localhost")
	vp.SetDefault("kafka.port", 9092)
	vp.SetDefault("kafka.channel", "test")

	// Set default Logger
	vp.SetDefault("logger.routes.events", "log_routes_events.log")
	vp.SetDefault("logger.routes.admin", "log_routes_admin.log")
	vp.SetDefault("logger.routes.auth", "log_routes_auth.log")
	vp.SetDefault("logger.routes.deck", "log_routes_deck.log")
	vp.SetDefault("logger.routes.game", "log_routes_game.log")
	vp.SetDefault("logger.routes.middleware", "log_routes_middleware.log")
	vp.SetDefault("logger.nft", "log_nft_service.log")
	vp.SetDefault("logger.blockchain", "log_blockchain_service.log")
	vp.SetDefault("logger.kafka", "log_kafka_service.log")
	vp.SetDefault("logger.matchmaking.service", "log_matchmaking_service.log")
	vp.SetDefault("logger.matchmaking.pool", "log_matchmaking_pool.log")
	vp.SetDefault("logger.match.casual", "log_match_casual.log")
	vp.SetDefault("logger.match.bot", "log_match_bot.log")
	vp.SetDefault("logger.match.private", "log_match_private.log")
	vp.SetDefault("logger.match.hardcore", "log_match_hardcore.log")
	vp.SetDefault("logger.match.manager", "log_match_matchmanager.log")
	vp.SetDefault("logger.match.service", "log_match_matchservice.log")
	vp.SetDefault("logger.websocket.topic.home", "log_websocket_topic_home.log")
	vp.SetDefault("logger.websocket.topic.match", "log_websocket_topic_match.log")
	vp.SetDefault("logger.server", "log_server.log")
	vp.SetDefault("cloud-config", false)

	//err = vp.ReadInConfig()
	//if err != nil {
	//	panic(err)
	//}

	if vp.GetBool("cloud-config") {

		err := vp.AddRemoteProvider("consul", "consul.cluster.internal:80", configKey)
		if err != nil {
			panic(err)
		}

		vp.SetConfigType("yaml")
		err = vp.ReadRemoteConfig()
		if err != nil {
			panic(err)
		}
	} else {
		err := vp.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}

	vp.AutomaticEnv()

	return vp
}
