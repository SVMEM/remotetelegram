package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"

	tb "gopkg.in/tucnak/telebot.v2"
)

type config struct {
	Token string `json:"token"`
	Admin int64  `json:"admin_id"`
}

func main() {
	cfg, err := openConfig()
	if err != nil {
		log.Fatalln(err)
	}
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatalln(err)
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	var n int = 0
	b.Handle("/start", func(m *tb.Message) {
		menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		menu.Reply(
			menu.Row(tb.Btn{Text: "Пробел или пауза"}),
			menu.Row(tb.Btn{Text: "Выключить"}),
			menu.Row(tb.Btn{Text: "MUTE"}),
			menu.Row(tb.Btn{Text: "-10sec"}),
			menu.Row(tb.Btn{Text: "+10sec"}),
			menu.Row(tb.Btn{Text: "-volume"}),
			menu.Row(tb.Btn{Text: "+volume"}))
		_, err := b.Send(m.Sender, "Главное меню", menu)
		if err != nil {
			log.Println(err)
		}
	})

	b.Handle("Пробел или пауза", func(m *tb.Message) {
		kb.SetKeys(keybd_event.VK_SPACE)
		if err := kb.Launching(); err != nil {
			log.Println(err)
		}
	})

	b.Handle("MUTE", func(m *tb.Message) {
		kb.SetKeys(keybd_event.VK_M)
		if err := kb.Launching(); err != nil {
			log.Println(err)
		}
	})

	b.Handle("+10sec", func(m *tb.Message) {
		kb.SetKeys(keybd_event.VK_RIGHT)
		if err := kb.Launching(); err != nil {
			log.Println(err)
		}
	})

	b.Handle("-10sec", func(m *tb.Message) {
		kb.SetKeys(keybd_event.VK_LEFT)
		if err := kb.Launching(); err != nil {
			log.Println(err)
		}
	})

	b.Handle("+volume", func(m *tb.Message) {
		switch runtime.GOOS {
		case "windows":
			cmd := exec.Command("rundll32.exe powrprof.dll, SetSuspendState Sleep")
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		case "darwin":
			n := n+10
			cmd := exec.Command("sudo", "osascript", "-e", "'set volume output volume", string(n), "'")
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		default:
			log.Fatalln("Error: unsupported OS.")
		}
	})

	b.Handle("-volume", func(m *tb.Message) {
		switch runtime.GOOS {
		case "windows":
			cmd := exec.Command("rundll32.exe powrprof.dll, SetSuspendState Sleep")
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		case "darwin":
			n := n-10
			cmd := exec.Command("sudo", "osascript", "-e", "'set volume output volume", string(n), "'")
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		default:
			log.Fatalln("Error: unsupported OS.")
		}
	})
	b.Handle("Выключить", func(m *tb.Message) {
		switch runtime.GOOS {
		case "windows":
			cmd := exec.Command("rundll32.exe powrprof.dll, SetSuspendState Sleep")
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		case "darwin":
			cmd := exec.Command("sudo", "shutdown", "-s", "now")
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		default:
			log.Fatalln("Error: unsupported OS.")
		}
	})

	b.Start()
}

func openConfig() (*config, error) {
	cfgFile, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	cfg := &config{}
	defer cfgFile.Close()
	err = json.NewDecoder(cfgFile).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func saveConfig(cfg *config) error {
	cfgData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile("config.json", cfgData, 0644)
}
