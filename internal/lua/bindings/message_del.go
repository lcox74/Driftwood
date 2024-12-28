package bindings

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	lua "github.com/yuin/gopher-lua"
)

// MessageBindingDelete provides Lua bindings for managing Discord messages.
type MessageBindingDelete struct {
	Session *discordgo.Session
}

// NewMessageBindingDelete initializes a new message management instance.
func NewMessageBindingDelete(session *discordgo.Session) *MessageBindingDelete {
	slog.Info("Creating new MessageBindingDel")
	return &MessageBindingDelete{
		Session: session,
	}
}

// Name returns the name of the binding.
func (b *MessageBindingDelete) Name() string {
	return "delete_message"
}

// Register registers the message-related functions in the Lua state.
func (b *MessageBindingDelete) Register(L *lua.LState) *lua.LFunction {
	return L.NewFunction(func(L *lua.LState) int {
		messageID := L.CheckString(1)
		channelID := L.CheckString(2)

		err := b.Session.ChannelMessageDelete(channelID, messageID)
		if err != nil {
			slog.Error("Failed to delete message", "message_id", messageID, "channel_id", channelID, "error", err)
			L.Push(lua.LString(fmt.Sprintf("Failed to delete message: %s", err.Error())))
			return 1
		}

		L.Push(lua.LString("Message deleted successfully"))
		return 1
	})
}

// HandleInteraction is not applicable for this binding.
func (b *MessageBindingDelete) HandleInteraction(L *lua.LState, interaction *discordgo.InteractionCreate) error {
	// This binding does not handle interactions
	return nil
}

func (b *MessageBindingDelete) CanHandleInteraction(interaction *discordgo.InteractionCreate) bool {
	return false
}