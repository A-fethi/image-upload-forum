package utils

import "fmt"

var (
	host    = "localhost"
	port    = 8080
	cyan    = "\033[1;36m"
	yellow  = "\033[1;33m"
	green   = "\033[1;32m"
	magenta = "\033[1;35m"
	reset   = "\033[0m"
	url     = fmt.Sprintf("http://%s:%d", host, port)
)

func Print() {
	fmt.Println(cyan + "╔════════════════════════════════════════════╗" + reset)
	fmt.Printf(cyan+"║ %-52s ║\n", yellow+"🚀 Server Information"+reset)
	fmt.Println(cyan + "╠════════════════════════════════════════════╣" + reset)
	fmt.Printf(cyan+"║ %-20s: %-37s ║\n"+reset, green+"🌐 URL", yellow+url+reset)
	fmt.Printf(cyan+"║ %-20s: %-37s ║\n"+reset, green+"⚡ Status", green+"Running"+reset)
	fmt.Printf(cyan+"║ %-20s: %-37s ║\n"+reset, green+"🔧 Environment", magenta+"Development"+reset)
	fmt.Println(cyan + "╠════════════════════════════════════════════╣" + reset)
	fmt.Printf(cyan+"║ %-42s║\n"+reset, "📄 Press CTRL+C to stop the server")
	fmt.Println(cyan + "╚════════════════════════════════════════════╝" + reset)
}
