package commands

import "github.com/bwmarrin/discordgo"

// OptInt reads an integer-valued interaction option, returning def when absent.
func OptInt(opts map[string]*discordgo.ApplicationCommandInteractionDataOption, name string, def int) int {
	if o, ok := opts[name]; ok && o != nil {
		if v, ok := o.Value.(float64); ok {
			return int(v)
		}
	}
	return def
}

// OptBool reads a boolean-valued interaction option.
func OptBool(opts map[string]*discordgo.ApplicationCommandInteractionDataOption, name string) bool {
	if o, ok := opts[name]; ok && o != nil {
		if v, ok := o.Value.(bool); ok {
			return v
		}
	}
	return false
}

// OptUser reads a user-valued interaction option, returning its ID.
func OptUser(opts map[string]*discordgo.ApplicationCommandInteractionDataOption, name string) string {
	if o, ok := opts[name]; ok && o != nil {
		if v, ok := o.Value.(string); ok {
			return parseMentionID(v)
		}
	}
	return ""
}
