package bot

// TODO: complete help function
func help() string {
	return "Commands:\n```" +
		"server start\nserver stop\n\n" +
		"azure login\nazure logout\nazure account\n\n" +
		"dj join\ndj leave\ndj play <song>\ndj add <song>\ndj insert <index> <song>\ndj remove <index>\ndj list\ndj skip\ndj clear\n\n" +
		"@ me to talk to me :D```"
}
