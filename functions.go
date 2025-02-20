package ettcodesdk

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"
	tgbotapiK "gopkg.in/telegram-bot-api.v4"
)

type ToolFunction struct {
	Cfg    *Config
	Logger *FaasLogger
}

func New(cfg *Config) *ToolFunction {
	return &ToolFunction{
		Cfg:    cfg,
		Logger: NewLoggerFunction(cfg.FunctionName),
	}
}

func (o *ToolFunction) SendTelegram(text string) error {
	client := &http.Client{}

	if ContainsLike(Mode, text) {
		text = strings.Replace(text, "\n", "", -1)
	} else {
		text = o.Cfg.FunctionName + " >>> " + time.Now().Format(time.RFC3339) + " >>>>> " + text
	}

	if o.Cfg.BranchName != "" {
		text = strings.ToUpper(o.Cfg.BranchName) + " >>> " + text
	}

	for _, e := range o.Cfg.AccountIds {
		botUrl := fmt.Sprintf("https://api.telegram.org/bot"+o.Cfg.BotToken+"/sendMessage?chat_id="+e+"&text=%s", text)
		request, err := http.NewRequest("GET", botUrl, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(request)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	return nil
}

func (o *ToolFunction) SendTelegramFile(req []byte, filename string) error {
	err := os.WriteFile(filename, req, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(filename)

	for _, e := range o.Cfg.AccountIds {
		bot, err := tgbotapiK.NewBotAPI(o.Cfg.BotToken)
		if err != nil {
			return err
		}

		message := tgbotapiK.NewDocumentUpload(cast.ToInt64(e), filename)
		_, err = bot.Send(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *ToolFunction) Config() *Config {
	return o.Cfg
}
