package config

type Jira struct {
	TaskBrowseURL string `env:"JIRA_TASK_BROWSE_URL" envDefault:""`
}

type Config struct {
	Telegram struct {
		Secret      string `env:"TELEGRAM_SECRET,notEmpty"`
		CallbackUrl string `env:"TELEGRAM_CALLBACK_URL" envDefault:""`
	}
	Jira Jira
}
