package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber"
)

type GHook struct {
	DashboardId int `json:"dashboardId"`
	EvalMatches []struct {
		Value  float64     `json:"value"`
		Metric string      `json:"metric"`
		Tags   interface{} `json:"tags"`
	} `json:"evalMatches"`
	ImageUrl string      `json:"imageUrl"`
	Message  string      `json:"message"`
	OrgId    int         `json:"orgId"`
	PanelId  int         `json:"panelId"`
	RuleId   int         `json:"ruleId"`
	RuleName string      `json:"ruleName"`
	RuleUrl  string      `json:"ruleUrl"`
	State    string      `json:"state"`
	Tags     interface{} `json:"tags"`
	Title    string      `json:"title"`
}

type AHook struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			AlertName    string `json:"alertName"`
			Alertname    string `json:"alertname"`
			Notification string `json:"notification"`
		} `json:"labels"`
		Annotations struct {
			AlertName    string `json:"alertName"`
			Content      string `json:"content"`
			Describe     string `json:"describe"`
			Notification string `json:"notification"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		AlertName    string `json:"alertName"`
		Alertname    string `json:"alertname"`
		Notification string `json:"notification"`
	} `json:"groupLabels"`
	CommonLabels struct {
		AlertName    string `json:"alertName"`
		Alertname    string `json:"alertname"`
		Notification string `json:"notification"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		AlertName    string `json:"alertName"`
		Content      string `json:"content"`
		Describe     string `json:"describe"`
		Notification string `json:"notification"`
	} `json:"commonAnnotations"`
	ExternalURL string `json:"externalURL"`
}

var sent_count int = 1
var color string

const (
	Url         = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	OKMsg       = "告警恢复"
	AlertingMsg = "触发告警"
	GfOK        = "OK"
	GfAlerting  = "Alerting"
	AlOK        = "resolved"
	AlAlerting  = "firing"
	ColorGreen  = "info"
	ColorGray   = "comment"
	ColorRed    = "warning"
)

func GwStat() func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		stat_msg := "G3WW Server created by Nova Kwok is running! \nParsed & forwarded " + strconv.Itoa(sent_count) + " messages to WeChat Work!"
		c.Send(stat_msg)
	}
}

func SedColor(title, ok, okmsg, alert, alertmsg string) (string, Color string) {
	if strings.Contains(title, ok) {
		title = strings.ReplaceAll(title, ok, okmsg)
		Color = ColorGreen
	} else {
		title = strings.ReplaceAll(title, alert, alertmsg)
		Color = ColorRed
	}
	return title, Color
}

func GwWorker() func(c *fiber.Ctx) {

	return func(c *fiber.Ctx) {
		some_app := c.Params("apps")
		url := Url + c.Params("key")
		var jsonStr []byte
		TIMEFORMATE := "2006-01-02 15:04:05"

		if some_app == "alertmanager" {
			var h *AHook = new(AHook)
			var OK = AlOK
			var Alerting = AlAlerting

			if err := c.BodyParser(h); err != nil {
				fmt.Println(err)
				c.Send("Error on JSON format")
				return
			}
			h.Status, color = SedColor(h.Status, OK, OKMsg, Alerting, AlertingMsg)

			if h.Alerts[0].Status == "resolved" {
				okMsgStr := fmt.Sprintf(`
			{
			  "msgtype": "markdown",
			  "markdown": {
				"content": "## <font color=\"%s\">%s</font>\n
				   >策略名称: <font color=\"comment\">%s</font>
				   >规则名称: **<font color=\"comment\">%s</font>**
				   >规则描述: **<font color=\"comment\">%s</font>**
				   >告警内容: **<font color=\"comment\">%s</font>**
				   >开始时间: <font color=\"comment\">%s</font>
				   >结束时间: <font color=\"comment\">%s</font>"
				}
			}
			`, color, h.Status, h.Alerts[0].Annotations.AlertName, h.Alerts[0].Labels.Alertname, h.Alerts[0].Annotations.Describe, h.Alerts[0].Annotations.Content, h.Alerts[0].StartsAt.Local().Format(TIMEFORMATE), h.Alerts[0].EndsAt.Local().Format(TIMEFORMATE))
				jsonStr = []byte(okMsgStr)

			} else if h.Alerts[0].Status == "firing" {
				alertMsgStr := fmt.Sprintf(`
			{
			  "msgtype": "markdown",
			  "markdown": {
				"content": "## <font color=\"%s\">%s</font>\n
				   >策略名称: <font color=\"comment\">%s</font>
				   >规则名称: **<font color=\"comment\">%s</font>**
				   >规则描述: **<font color=\"comment\">%s</font>**
				   >告警内容: **<font color=\"comment\">%s</font>**
				   >开始时间: <font color=\"comment\">%s</font>"
				}
			}
			`, color, h.Status, h.Alerts[0].Annotations.AlertName, h.Alerts[0].Labels.Alertname, h.Alerts[0].Annotations.Describe, h.Alerts[0].Annotations.Content, h.Alerts[0].StartsAt.Local().Format(TIMEFORMATE))
				jsonStr = []byte(alertMsgStr)
			} else {
				fmt.Println("The alert.status didn't match, please check.")
			}

		} else if some_app == "grafana" {
			var h *GHook = new(GHook)
			var OK = GfOK
			var Alerting = GfAlerting

			NewRuleUrl := strings.Replace(h.RuleUrl, "http://localhost:3000", "https://yourown.com", -1)

			h.Title, color = SedColor(h.Title, OK, OKMsg, Alerting, AlertingMsg)

			if err := c.BodyParser(h); err != nil {
				fmt.Println(err)
				c.Send("Error on JSON format")
				return
			}

			if h.State == "ok" {
				okMsgStr := fmt.Sprintf(`
		{
			"msgtype": "markdown",
			"markdown": {
			  "content": "<font color=\"%s\">%s</font>\r\n<font color=\"comment\">%s\r\n[点击查看详情](%s)</font>"
			}
		  }
		`, color, h.Title, h.Message, NewRuleUrl)

				jsonStr = []byte(okMsgStr)

			} else {
				alertMsgStr := fmt.Sprintf(`
		{
			"msgtype": "markdown",
			"markdown": {
			  "content": "<font color=\"%s\">%s</font>\r\n<font color=\"comment\">%s\r\n报警数值为：%f\r\n[点击查看详情](%s)</font>"
			}
		  }
		`, color, h.Title, h.Message, h.EvalMatches[0].Value, NewRuleUrl)
				jsonStr = []byte(alertMsgStr)
			}

		} else {
			fmt.Printf("The urp-path apps didn't match any of apps vars.")
			return
		}

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.Send("Error sending to WeChat Work API")
			return
		}
		defer resp.Body.Close()
		c.Send(resp)
		sent_count++
	}
}
