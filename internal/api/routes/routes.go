package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/peithosecure/peitho-backend/internal/api/handlers"
	"github.com/peithosecure/peitho-backend/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health check endpoints
	r.HandleFunc("/healthz", handlers.HealthzHandler).Methods(http.MethodGet)
	r.HandleFunc("/status", handlers.StatusHandler).Methods(http.MethodGet)

	// Auth routes (including unlock endpoints)
	authRouter := r.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.HandleFunc("/register", handlers.RegisterHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/refresh", handlers.RefreshHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/logout", handlers.LogoutHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/delete", handlers.DeleteAccountHandler).Methods(http.MethodDelete)
	authRouter.HandleFunc("/check", handlers.CheckEmailHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/verify-email", handlers.VerifyEmailHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/resend-token", handlers.ResendVerificationTokenHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/send-verification", handlers.SendVerificationLinkHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/request-password-reset", handlers.RequestPasswordResetHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/reset-password", handlers.ResetPasswordHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/setup-password", handlers.SetupPasswordHandler).Methods(http.MethodPost)

	// Add unlock-status and unlock-validate here inside authRouter
	authRouter.HandleFunc("/unlock-status", handlers.UnlockStatusHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/unlock/validate", handlers.UnlockValidateHandler).Methods(http.MethodPost)

	// Secure routes - protected by AuthGuard middleware
	secureRouter := r.PathPrefix("/api/v1/auth/secure-sample").Subrouter()
	secureRouter.Use(middleware.AuthGuard)
	secureRouter.HandleFunc("", handlers.SecureSampleHandler).Methods(http.MethodGet)

	// PQC locked routes - group with UnlockGuardMiddleware for cleaner code
	pqcRouter := r.PathPrefix("/api/v1").Subrouter()
	pqcRouter.Use(middleware.UnlockGuardMiddleware)
	pqcRouter.HandleFunc("/events/log", handlers.EngineEventHandler).Methods(http.MethodPost)
	pqcRouter.HandleFunc("/security-scan", handlers.ProwlerScanHandler).Methods(http.MethodGet)
	pqcRouter.HandleFunc("/metrics", handlers.MetricsHandler).Methods(http.MethodGet)
	pqcRouter.HandleFunc("/admin-metrics", handlers.AdminMetricsHandler).Methods(http.MethodGet)

	// Token metrics (public API)
	r.HandleFunc("/api/v1/metrics/tokens", handlers.TokenMetricsHandler).Methods(http.MethodGet)

	// App integrations (Bearer-protected)
	r.Handle("/api/v1/integrations",
		middleware.AuthGuard(http.HandlerFunc(handlers.GetAppIntegrations)),
	).Methods(http.MethodGet)

	// Audit analytics (Bearer-protected)
	r.Handle("/api/v1/analytics/audit",
		middleware.AuthGuard(http.HandlerFunc(handlers.AuditAnalyticsHandler)),
	).Methods(http.MethodGet)

	// Deep links - **Note:** two handlers for same path and method will conflict
	// Consider merging or using distinct paths or HTTP methods
	r.HandleFunc("/api/v1/deeplink", handlers.DeepLinkRedirectHandler).Methods(http.MethodGet)
	// If UniversalDeeplinkHandler differs, consider different path or method
	// r.HandleFunc("/api/v1/deeplink/universal", handlers.UniversalDeeplinkHandler).Methods(http.MethodGet)

	// Trace log (Bearer + PQC Unlock)
	r.Handle("/api/v1/log/trace",
		middleware.UnlockGuardMiddleware(
			middleware.AuthGuard(http.HandlerFunc(handlers.TraceLogHandler))),
	).Methods(http.MethodGet)

	// Swagger Docs endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Log all registered routes (handle error)
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		fmt.Printf("📡 Registered route: %s [%v]\n", path, methods)
		return nil
	})
	if err != nil {
		log.Printf("⚠️ Error walking routes: %v", err)
	}

	return r
}
