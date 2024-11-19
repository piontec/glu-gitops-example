package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var SHA string

func main() {
	body := fmt.Sprintf(`<html>
    <body>
        <h1>Welcome!</h1>
        <p><strong>Environment: %q</strong></p>
        <p>Build SHA: %q</p>
        <p>Image Digest: %q</p>
    </body>
</html>`,
		os.Getenv("APP_ENVIRONMENT"),
		SHA,
		os.Getenv("APP_IMAGE_DIGEST"),
	)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}))

	slog.Info("Starting server", "addr", ":8080")

	http.ListenAndServe(":8080", nil)
}
