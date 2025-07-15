package commands

import (
	"fmt"
	"luna/i18n"
	"luna/interfaces"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct {
	StartTime time.Time
	Store     interfaces.DataStore
}

func (c *PingCommand) GetCommandDef() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "BOTの応答速度や状態を確認します",
	}
}

func (c *PingCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lang := i.Locale

	start := time.Now()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: i18n.GetMessage(lang, "ping_command.pinging", nil),
		},
	})
	if err != nil {
		return
	}
	apiLatency := time.Since(start)

	wsLatency := s.HeartbeatLatency()

	dbStart := time.Now()
	dbErr := c.Store.PingDB()
	dbLatency := time.Since(dbStart)
	dbStatusKey := "ping_command.db_online"
	if dbErr != nil {
		dbStatusKey = "ping_command.db_offline"
	}
	dbStatus := i18n.GetMessage(lang, dbStatusKey, nil)

	uptime := time.Since(c.StartTime)

	embed := &discordgo.MessageEmbed{
		Title: i18n.GetMessage(lang, "ping_command.title", nil),
		Color: 0x7289da, // Discord Blue
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   i18n.GetMessage(lang, "ping_command.api_latency", nil),
				Value:  fmt.Sprintf("```%s```", apiLatency.Round(time.Millisecond).String()),
				Inline: true,
			},
			{
				Name:   i18n.GetMessage(lang, "ping_command.ws_latency", nil),
				Value:  fmt.Sprintf("```%s```", wsLatency.Round(time.Millisecond).String()),
				Inline: true,
			},
			{
				Name:   i18n.GetMessage(lang, "ping_command.database", nil),
				Value:  fmt.Sprintf("```%s (%s)```", dbStatus, dbLatency.Round(time.Millisecond).String()),
				Inline: true,
			},
			{
				Name:  i18n.GetMessage(lang, "ping_command.uptime", nil),
				Value: fmt.Sprintf("```%s```", formatUptime(uptime)),
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: new(string),
		Embeds:  &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		// c.Log is not available in PingCommand.
	}
}

// アップタイムを見やすい形式にフォーマットするヘルパー関数
func formatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d日", days))
	}
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%d時間", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%d分", m))
	}
	if s > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d秒", s))
	}

	return strings.Join(parts, " ")
}

func (c *PingCommand) HandleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {}
func (c *PingCommand) HandleModal(s *discordgo.Session, i *discordgo.InteractionCreate)     {}
func (c *PingCommand) GetComponentIDs() []string                                            { return []string{} }
func (c *PingCommand) GetCategory() string                                                  { return "ユーティリティ" }
