package vocman
// vocman - VOCabulary MANagment

func appendMessageToFile(filename, message string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if strings.Contains(message, "http://") || strings.Contains(message, "https://") || strings.Contains(message, "@") {
		log.Println("Message contains link or username! Skipping.")
		return
	} else {}

	if _, err := file.WriteString(message + "\n"); err != nil {
		log.Println("Error writing to file:", err)
	}
}