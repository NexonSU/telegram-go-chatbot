package utils

import "os"

type Configuration struct {
	Telegram struct {
		Token         string   `json:"token,omitempty"`
		Chat          string   `json:"chat,omitempty"`
		StreamChannel string   `json:"stream_channel,omitempty"`
		Channel       string   `json:"channel,omitempty"`
		BotApiUrl     string   `json:"bot_api_url,omitempty"`
		Admins        []string `json:"admins,omitempty"`
		Moders        []string `json:"moders,omitempty"`
		SysAdmin      string   `json:"sysadmin,omitempty"`
	}
	Webhook struct {
		Listen         string   `json:"listen,omitempty"`
		Port           int      `json:"port,omitempty"`
		AllowedUpdates []string `json:"allowed_updates,omitempty"`
	}
	Youtube struct {
		ApiKey      string `json:"api_key,omitempty"`
		ChannelName string `json:"channel_name,omitempty"`
		ChannelID   string `json:"channel_id,omitempty"`
	}
	CurrencyKey string `json:"currency_key,omitempty"`
	ReleasesUrl string `json:"releases_url,omitempty"`
}

var ConfigFile, _ = os.Open("../config.json")
var Config = new(Configuration)
