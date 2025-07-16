package commands

import (
	"fmt"
	"luna/i18n"
	"luna/interfaces"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type UserInfoCommand struct {
	Log interfaces.Logger
}

func (c *UserInfoCommand) GetCommandDef() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "user-info",
		Description: "指定したユーザーの情報を表示します",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "情報を表示するユーザー", Required: false},
		},
	}
}

func (c *UserInfoCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lang := i.Locale
	options := i.ApplicationCommandData().Options
	var targetUser *discordgo.User
	if len(options) > 0 {
		targetUser = options[0].UserValue(s)
	} else {
		targetUser = i.Member.User
	}

	member, err := s.State.Member(i.GuildID, targetUser.ID)
	if err != nil {
		member, err = s.GuildMember(i.GuildID, targetUser.ID)
		if err != nil {
			c.Log.Error("メンバー情報の取得に失敗", "error", err, "userID", targetUser.ID)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: i18n.GetMessage(lang, "user_info_command.error_fetch", nil), Flags: discordgo.MessageFlagsEphemeral},
			})
			return
		}
	}

	joinedAt := member.JoinedAt
	createdAt, _ := discordgo.SnowflakeTimestamp(targetUser.ID)

	roles := make([]string, 0)
	guildRoles, _ := s.GuildRoles(i.GuildID)
	for _, roleID := range member.Roles {
		for _, role := range guildRoles {
			if role.ID == roleID {
				roles = append(roles, fmt.Sprintf("<@&%s>", role.ID))
				break
			}
		}
	}
	rolesStr := i18n.GetMessage(lang, "user_info_command.none", nil)
	if len(roles) > 0 {
		rolesStr = strings.Join(roles, " ")
	}

	presence, err := s.State.Presence(i.GuildID, targetUser.ID)
	statusStr := i18n.GetMessage(lang, "user_info_command.status_offline", nil)
	activityStr := i18n.GetMessage(lang, "user_info_command.none", nil)
	if err == nil {
		statusStr = i18n.GetMessage(lang, "user_info_command.status_"+string(presence.Status), nil)

		if len(presence.Activities) > 0 {
			activity := presence.Activities[0]
			activityType := i18n.GetMessage(lang, "user_info_command.activity_"+strings.ToLower(activity.Type.String()), nil)
			activityStr = fmt.Sprintf("%s: %s", activityType, activity.Name)
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:     i18n.GetMessage(lang, "user_info_command.title", map[string]interface{}{"Username": targetUser.Username}),
		Color:     s.State.UserColor(targetUser.ID, i.ChannelID),
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: member.AvatarURL("1024")},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    targetUser.String(),
			IconURL: targetUser.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: i18n.GetMessage(lang, "user_info_command.section_basic", nil), Value: fmt.Sprintf("**%s** `%s`\n**%s** %v", i18n.GetMessage(lang, "user_info_command.field_id", nil), targetUser.ID, i18n.GetMessage(lang, "user_info_command.field_bot", nil), targetUser.Bot), Inline: false},
			{Name: i18n.GetMessage(lang, "user_info_command.section_dates", nil), Value: fmt.Sprintf("**%s** <t:%d:R>\n**%s** <t:%d:R>", i18n.GetMessage(lang, "user_info_command.field_account_created", nil), createdAt.Unix(), i18n.GetMessage(lang, "user_info_command.field_server_joined", nil), joinedAt.Unix()), Inline: false},
			{Name: i18n.GetMessage(lang, "user_info_command.section_status", nil), Value: statusStr, Inline: true},
			{Name: i18n.GetMessage(lang, "user_info_command.section_activity", nil), Value: activityStr, Inline: true},
			{Name: i18n.GetMessage(lang, "user_info_command.section_roles", map[string]interface{}{"Count": len(roles)}), Value: rolesStr, Inline: false},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

func (c *UserInfoCommand) HandleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {}
func (c *UserInfoCommand) HandleModal(s *discordgo.Session, i *discordgo.InteractionCreate)     {}
func (c *UserInfoCommand) GetComponentIDs() []string                                            { return []string{} }
func (c *UserInfoCommand) GetCategory() string {
	return "ユーティリティ"
}
