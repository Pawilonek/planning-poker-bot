# Scrumpoke

Add Telegram bot: [Scrumpoke](https://t.me/ScrumpokeBot)

### Local development:
1. Run ngrok (https://ngrok.com/) with command and get generated address:
`ngrok http 8080`
2. Setup your bot token environment variable:
`export TELEGRAM_SECRET=123:Secret!!!`
3. Set your callback URL:
`export TELEGRAM_CALLBACK_URL=https://randomUUID.ngrok.io`
4. Setup Callback URL for your bot
`curl -X POST "https://api.telegram.org/bot$TELEGRAM_SECRET/setWebhook?url=$TELEGRAM_CALLBACK_URL/webhook"`
5. Run bot:
`./bin/scrumpoke`
