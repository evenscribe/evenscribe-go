package main

type LogOwner struct {
	HostName string `json:"host_name"` /* Identifier for frontend */
	AppName  string `json:"app_name"`  /* Identifier for backend */
}

type Log struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type LogEntry struct {
	Timestamp string   `json:"@timestamp"`
	LogOwner  LogOwner `json:"log_owner"`
	Log       Log      `json:"log"`
	Message_  string   `json:"_msg"`
}
