package handler

import (
	"log"
	"net/http"
	"smply/config"
	"smply/internal/queue"
	"smply/internal/render"
	"smply/internal/service"
	"strings"
)

func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	render.SinglePage(w, "admin/login.html", render.ViewData{})
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	if email == "" {
		render.SinglePage(w, "admin/login.html", render.ViewData{
			Error: "Email is required",
		})
		return
	}

	token, err := service.CreateMagicToken(r.Context(), email)
	if err != nil {
		log.Printf("Error creating magic token for %s: %v", email, err)
		render.SinglePage(w, "admin/login.html", render.ViewData{
			Error: "Internal server error",
		})
		return
	}

	err = queue.EnqueueAdminLoginMagicLinkEmail(r.Context(), email, token)
	if err != nil {
		log.Println("Error enqueuing admin login magic link email:", err)
		render.SinglePage(w, "admin/login.html", render.ViewData{
			Error: "Internal server error",
		})
		return
	}

	render.SinglePage(w, "admin/login.html", render.ViewData{
		Data: map[string]string{
			"SentTo": email,
		},
	})
}

func AdminAuth(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		render.ErrorPage(w, http.StatusUnprocessableEntity, "Cannot find token")
		return
	}

	magicToken, err := service.ValidateMagicToken(r.Context(), token)
	if err != nil {
		log.Printf("Error validating magic token: %v", err)
		render.ErrorPage(w, http.StatusUnprocessableEntity, "Invalid or expired token")
		return
	}

	sessionId, err := service.CreateSession(r.Context(), magicToken.Email)
	if err != nil {
		log.Println("Error creating session:", err)
		render.SinglePage(w, "admin/login.html", render.ViewData{
			Error: "Internal server error",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     strings.ToLower(config.Env.App.Name) + "_session_id",
		Value:    sessionId,
		Path:     "/",
		MaxAge:   60 * 60,
		Secure:   true,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/admin/stats", http.StatusFound)
}

func AdminStats(w http.ResponseWriter, r *http.Request) {
	render.SinglePage(w, "admin/stats.html", render.ViewData{})
	// stats, err := service.GetAdminStats(r.Context())
	// if err != nil {
	// 	log.Println(err)
	// 	render.ErrorPage(w, http.StatusInternalServerError, "Internal server error")
	// 	return
	// }

	// t := template.New("").Funcs(template.FuncMap{
	// 	"add": func(a, b int) int { return a + b },
	// 	"percent": func(val, max int64) int64 {
	// 		if max == 0 {
	// 			return 0
	// 		}
	// 		return val * 100 / max
	// 	},
	// 	"lower": strings.ToLower,
	// })

	// t, err = t.ParseFiles("templates/pages/admin-stats.html")
	// if err != nil {
	// 	log.Println(err)
	// 	render.ErrorPage(w, http.StatusInternalServerError, "Internal server error")
	// 	return
	// }

	// if err := t.ExecuteTemplate(w, "admin-stats.html", render.ViewData{
	// 	Title: "Admin Stats",
	// 	Page:  "admin_stats",
	// 	Data: map[string]any{
	// 		"Stats": stats,
	// 	},
	// }); err != nil {
	// 	log.Println("EXECUTE ERROR:", err)
	// 	render.ErrorPage(w, http.StatusInternalServerError, "Internal server error")
	// 	return
	// }

	// render.Page(w, "admin/stats.html", render.ViewData{
	// Title: "Admin Stats",
	// Page:  "admin_stats",
	// Data: map[string]any{
	// 	"Stats": stats,
	// },
	// })
}
