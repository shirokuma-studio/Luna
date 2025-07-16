package commands

import (
	"fmt"
	"luna/i18n"

	"github.com/bwmarrin/discordgo"
)

type AvatarCommand struct{}

func (c *AvatarCommand) GetCommandDef() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "avatar",
		Description: "ユーザーのアバターやバナーを表示します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "情報を表示するユーザー",
				Required:    false,
			},
		},
	}
}

func (c *AvatarCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lang := i.Locale
	options := i.ApplicationCommandData().Options
	var targetUser *discordgo.User
	var targetMember *discordgo.Member

	if len(options) > 0 {
		targetUser = options[0].UserValue(s)
		m, err := s.State.Member(i.GuildID, targetUser.ID)
		if err != nil {
			m, err = s.GuildMember(i.GuildID, targetUser.ID)
			if err != nil {
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: i18n.GetMessage(lang, "avatar_command.error_fetch", nil), Flags: discordgo.MessageFlagsEphemeral},
				}); err != nil {
					fmt.Printf("Failed to respond to interaction: %v\n", err)
				}
				return
			}
		}
		targetMember = m
	} else {
		targetUser = i.Member.User
		targetMember = i.Member
	}

	userWithBanner, err := s.User(targetUser.ID)
	if err != nil {
		userWithBanner = targetUser
	}

	avatarURL := targetUser.AvatarURL("1024")
	serverAvatarURL := targetMember.AvatarURL("1024")
	bannerURL := userWithBanner.BannerURL("1024")

	embed := &discordgo.MessageEmbed{
		Title: i18n.GetMessage(lang, "avatar_command.title", map[string]interface{}{"Username": targetUser.Username}),
		Color: 0x7289da,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    targetUser.String(),
			IconURL: avatarURL,
		},
		Fields: []*discordgo.MessageEmbedField{},
	}

	embed.Image = &discordgo.MessageEmbedImage{URL: avatarURL}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: i18n.GetMessage(lang, "avatar_command.field_global_avatar", nil), Value: fmt.Sprintf("[%s](%s)", i18n.GetMessage(lang, "avatar_command.link", nil), avatarURL)})

	if targetMember.Avatar != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: serverAvatarURL}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: i18n.GetMessage(lang, "avatar_command.field_server_avatar", nil), Value: fmt.Sprintf("[%s](%s)", i18n.GetMessage(lang, "avatar_command.link", nil), serverAvatarURL)})
	}

	if userWithBanner.Banner != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: i18n.GetMessage(lang, "avatar_command.field_banner", nil), Value: fmt.Sprintf("[%s](%s)", i18n.GetMessage(lang, "avatar_command.link", nil), bannerURL)})
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}); err != nil {
		fmt.Printf("Failed to respond to interaction: %v\n", err)
	}
}

func (c *AvatarCommand) HandleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {}
func (c *AvatarCommand) HandleModal(s *discordgo.Session, i *discordgo.InteractionCreate)     {}
func (c *AvatarCommand) GetComponentIDs() []string                                            { return []string{} }
func (c *AvatarCommand) GetCategory() string                                                  { return "ユーティリティ" }
