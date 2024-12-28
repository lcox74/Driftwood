package utils

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	lua "github.com/yuin/gopher-lua"
)

// replyFunction returns a Lua function for replying to interactions.
// This utility can be used across multiple bindings.
func ReplyFunction(session *discordgo.Session, interaction *discordgo.InteractionCreate) lua.LGFunction {
	return func(L *lua.LState) int {
		argCount := L.GetTop()
		var message string
		var options *lua.LTable

		if argCount == 3 {
			L.CheckType(1, lua.LTTable) // Check 'self' argument is a table
			message = L.CheckString(2)
			L.CheckType(3, lua.LTTable) // Check 'options' argument is a table
			options = L.OptTable(3, nil)
		} else if argCount == 2 {
			L.CheckType(1, lua.LTTable) // Check 'self' argument is a table
			message = L.CheckString(2)
		} else {
			L.ArgError(1, "invalid arguments, expected (message [, options])")
			return 0
		}

		ephemeral := false
		mention := true
		if options != nil {
			if options.RawGetString("ephemeral") != lua.LNil {
				if options.RawGetString("ephemeral").Type() != lua.LTBool {
					L.ArgError(1, "'ephemeral' in options must be a boolean")
					return 0
				}
				ephemeral = lua.LVAsBool(options.RawGetString("ephemeral"))
			}
			if options.RawGetString("mention") != lua.LNil {
				if options.RawGetString("mention").Type() != lua.LTBool {
					L.ArgError(1, "'mention' in options must be a boolean")
					return 0
				}
				mention = lua.LVAsBool(options.RawGetString("mention"))
			}
		}

		if mention {
			message = fmt.Sprintf("<@%s> %s", interaction.Member.User.ID, message)
		}

		flags := discordgo.MessageFlags(0)
		if ephemeral {
			flags = discordgo.MessageFlagsEphemeral
		}

		if err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
				Flags:   flags,
			},
		}); err != nil {
			slog.Error("Failed to send interaction reply", "error", err)
		}

		return 0
	}
}