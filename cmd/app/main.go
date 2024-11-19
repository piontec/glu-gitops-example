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
    <head>
    <style>
        h1 {
            margin: 0 0 1rem 0;
        }

        p {
            margin: 4px 0;
        }

        body {
            display: flex;
            background-color: #62979D;
        }

        .container {
            display: flex;
            flex-direction: column;
            margin: 0 auto;
            border-radius: 1rem;
            padding: 1rem;
            background-color: #82ACB0;
            color: #9D6862;
        } 
    </style>
    </head>
    <body>
        <div class="container">
            <h1>Welcome!</h1>
            <p><strong>Environment: %q</strong></p>
            <p>Build SHA: %q</p>
            <p>Image Digest: %q</p>
        </div>
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
