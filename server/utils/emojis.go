package utils

var emojiMap = map[string]string{
	"listenning": "📡 ",
	"info":       "ℹ️ ",
	"connection": "🎯 ",
	"closed":     "🚪 ",
	"not ok":     "❌ ",
	"ok":         "✔️ ",
	"setting":    "⚙️ ",
	"debug":      "🛠️ ",
	"task":       "🎯 ",
	"alert":      "🔔 ",
	"loading":    "⏳ ",
	"send":       "🚀 ",
	"user":       "🤖 ",
	"error":      "❗ ",
	"help":       "💡 ",
	"input":      "💬 ",
}

func GetEmoji(key string) string {
	if emoji, ok := emojiMap[key]; ok {
		return emoji
	}
	return "ℹ️" // Emoji par défaut si la clé n'existe pas
}
