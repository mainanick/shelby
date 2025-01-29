package constants

const (
	PATH = "/etc/shellby"

	TRAEFIK_PATH         = PATH + "/traefik"
	DYNAMIC_TRAEFIK_PATH = TRAEFIK_PATH + "/dynamic"

	LOG_PATH = PATH + "/log"
	SSH_PATH = PATH + "/ssh"

	// Traefik
	TRAEFIK_FILE         = TRAEFIK_PATH + "/traefik.yml"
	SHELLBY_TRAEFIK_FILE = DYNAMIC_TRAEFIK_PATH + "/shellby.yml"

	// Log
	LOG_FILE = LOG_PATH + "/shellby.log"
)
