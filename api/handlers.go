// Package api provides HTTP routing and middleware for the authserver-backend application.
// It defines the routes and handlers for the authentication server.
// It also includes middleware for CORS and authentication.
// The handlers interact with the database via the dbrepo package and manage JWT tokens via the auth package.

package api

import (
	"authserver-backend/auth"
	"authserver-backend/internal/dbrepo"
	"authserver-backend/internal/models"
	"authserver-backend/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
)

// AuthServerApp holds the application configuration and dependencies
// such as the database repository, authentication service, and JWT settings.
// It is passed to handler functions to access these shared resources.
// Adjust the fields as necessary to fit your application's needs.
// Make sure to initialize this struct properly in your main application setup.
type AuthServerApp struct {
	DSN          string
	Domain       string
	DB           dbrepo.DatabaseRepo
	User         models.User
	Auth         auth.Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

// Home is a simple handler that responds with a JSON payload indicating the service status.
// It can be used to verify that the server is running and reachable.

func (app *AuthServerApp) Home(w http.ResponseWriter, r *http.Request) {

	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go apps up and running",
		Version: "1.0.0",
	}

	_ = utils.JSONResponse.WriteJSON(utils.JSONResponse{}, w, http.StatusOK, payload)
}

// Apps retrieves all apps from the database and returns them as a JSON response.
// If an error occurs during the database query, it responds with an error message and appropriate HTTP status code.
// This handler is not publicly accessible and does require user authentication.
// It is intended for internal use or for authenticated users only.
// Consider adding authentication middleware to protect this endpoint.
// Note: This function is similar to AppsCatalogue but is intended for non-admin use.
// This probably will give access to a limited set of apps based on user roles in the future and
// from here, to the links to the actual apps.
func (app *AuthServerApp) Apps(w http.ResponseWriter, r *http.Request) {

	apps, err := app.DB.AllApps("")
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	_ = utils.JSONResponse.WriteJSON(utils.JSONResponse{}, w, http.StatusOK, apps)
}

// AppsCatalogue retrieves all apps from the database for administrative purposes and returns them as a JSON response.
// This handler is protected by authentication middleware to ensure only authorized admin can access it.
// If an error occurs during the database query, it responds with an error message.
// This is for admin use only.
func (app *AuthServerApp) AppsCatalogue(w http.ResponseWriter, r *http.Request) {

	apps, err := app.DB.AllApps("")
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}
	_ = utils.JSONResponse{}.WriteJSON(w, http.StatusOK, apps)

}

// GetApp retrieves a specific app by its ID from the database and returns it as a JSON response.
// The app ID is extracted from the URL parameters.
// If the ID is missing or invalid, or if an error occurs during the database query,
// it responds with an appropriate error message and HTTP status code.
// This handler is for common users to get app details and probably links to the actual app.
func (app *AuthServerApp) GetApp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("id is missing in URL"), http.StatusBadRequest)
		return
	}

	appID, err := strconv.Atoi(id)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	thisapp, err := app.DB.ThisApp(appID, "")
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	_ = utils.JSONResponse{}.WriteJSON(w, http.StatusOK, thisapp)
}

// ThisApp is similar to GetApp but is intended for internal use or different access levels.
// It is necessary refine its purpose and usage in the application context.
func (app *AuthServerApp) ThisApp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		// If the ID is missing / invalid, return a 400 Bad Request error
		utils.JSONResponse{}.ErrorJSON(w, errors.New("id is missing in the URL"), http.StatusBadRequest)
		return
	}

	appID, err := strconv.Atoi(id)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	thisapp, err := app.DB.ThisApp(appID, "")
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	_ = utils.JSONResponse{}.WriteJSON(w, http.StatusOK, thisapp)
}

// ThisAppForEdit retrieves a specific app by its ID for editing purposes.
// It is almost identical to GetApp but is intended for use in an admin context
// where the app details can be modified and for adding a new app.
func (app *AuthServerApp) ThisAppForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("id is missing in URL"), http.StatusBadRequest)
		return
	}

	appID, err := strconv.Atoi(id)
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	thisapp, err := app.DB.ThisApp(appID, "")
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	_ = utils.JSONResponse{}.WriteJSON(w, http.StatusOK, thisapp)
}

// InsertApp reads a new app's details from the request body, validates the input,
// and inserts the app into the database. It responds with a success message and the new app ID
// or an error message if the input is invalid or the database operation fails.
// This handler is intended for admin use only.
func (app *AuthServerApp) InsertApp(w http.ResponseWriter, r *http.Request) {
	var thisapp models.ThisApp

	err := utils.JSONResponse{}.ReadJSON(w, r, &thisapp)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}
	var newapp models.NewApp

	newapp.Name = thisapp.Name
	newapp.Release = thisapp.Release
	newapp.Path = thisapp.Path
	newapp.Init = thisapp.Init
	newapp.Web = thisapp.Web
	newapp.Title = thisapp.Title
	newapp.Created = time.Now().Unix()
	newapp.Updated = time.Now().Unix()

	newID, err := app.DB.InsertApp(newapp, "")
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{
		Error:   false,
		Message: "app inserted " + strconv.Itoa(newID),
	}
	utils.JSONResponse{}.WriteJSON(w, http.StatusAccepted, resp)
}

// UpdateApp reads an app's updated details from the request body, validates the input,
// and updates the app in the database. It responds with a success message or an error message
// if the input is invalid or the database operation fails.
// This handler is intended for admin use only.
func (app *AuthServerApp) UpdateApp(w http.ResponseWriter, r *http.Request) {
	var payload models.ThisApp

	err := utils.JSONResponse{}.ReadJSON(w, r, &payload)

	// fmt.Printf("%+v\n", payload)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	thisapp, err := app.DB.ThisApp(payload.ID, "")

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	thisapp.ID = payload.ID
	thisapp.Name = payload.Name
	thisapp.Release = payload.Release
	thisapp.Path = payload.Path
	thisapp.Init = payload.Init
	thisapp.Web = payload.Web
	thisapp.Title = payload.Title
	thisapp.Created = payload.Created
	thisapp.Updated = time.Now().Unix()

	err = app.DB.UpdateApp(*thisapp, "")

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{
		Error:   false,
		Message: "app updated: ",
	}
	utils.JSONResponse{}.WriteJSON(w, http.StatusAccepted, resp)
}

// DeleteApp deletes an app from the database based on the ID provided in the URL parameters.
// It responds with a success message or an error message if the ID is missing, invalid,
// or if the database operation fails.
func (app *AuthServerApp) DeleteApp(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("id is missing in URL"), http.StatusBadRequest)
		return
	}

	appID, err := strconv.Atoi(id)
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	err = app.DB.DeleteApp(appID)
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{
		Error:   false,
		Message: "app deleted",
	}
	utils.JSONResponse{}.WriteJSON(w, http.StatusAccepted, resp)
}

// GetReleases returns a list of available release options
func (app *AuthServerApp) GetReleases(w http.ResponseWriter, r *http.Request) {
	// Optional: Check for JWT token if releases are admin-only
	// (Similar to your other handlers, e.g., using middleware)

	// Release options with id equal to value (no conversion needed)
	// This is hardcoded for now in dbrepo, but could be fetched from a database table in the future
	// or from a configuration file.
	resp, err := app.DB.GetReleases()
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Respond with JSON
	_ = utils.JSONResponse.WriteJSON(utils.JSONResponse{}, w, http.StatusOK, resp)
}

// The rest of the functions are related to authentication and session management.
// They handle user login, token refresh, session validation, and logout.
// These handlers interact with the auth package to manage JWT tokens and cookies.
// Ensure that the auth package is properly configured and integrated with your user management system.
// It should be considered to move these functions to a separate file for better organization.

func (app *AuthServerApp) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload models.User

	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("user not found or database error"), http.StatusBadRequest)
		return
	}
	//check password
	valid, err := app.User.PasswordMatches(requestPayload.Password, user.Password)
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	} else if !valid {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("invalid password"), http.StatusUnauthorized)
		return
	}

	// create a jwt user
	u := auth.JWTUser{
		ID:    user.ID,
		Email: user.Email,
	}

	//generate tokens
	tokens, err := app.Auth.GenerateTokenPair(&u)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// set the refresh token in an http only cookie
	refreshCookie := app.Auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	utils.JSONResponse{}.WriteJSON(w, http.StatusAccepted, tokens)

}

// RefreshToken handles the refresh token process.
// It checks for a valid refresh token in the cookies, verifies it,
// and issues a new pair of access and refresh tokens if valid.
// The new refresh token is set in an HTTP-only cookie.
// If the refresh token is missing, invalid, or expired, it responds with an error message.
func (app *AuthServerApp) RefreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.Auth.CookieName {
			claims := &auth.Claims{}
			refreshToken := cookie.Value
			// check if the token is empty
			if refreshToken == "" {
				utils.JSONResponse{}.ErrorJSON(w, errors.New("no refresh token found"), http.StatusUnauthorized)
				return
			}
			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
				return []byte(app.JWTSecret), nil
			})

			if err != nil {
				utils.JSONResponse{}.ErrorJSON(w, errors.New("no secret was found"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims

			ID := claims.UserID
			user, err := app.DB.GetUserByID(ID)
			if err != nil {
				http.Error(w, "Unknown user", http.StatusUnauthorized)
				return
			}
			u := auth.JWTUser{
				ID:    user.ID,
				Email: user.Email,
			}

			tokenPairs, err := app.Auth.GenerateTokenPair(&u)
			if err != nil {
				utils.JSONResponse{}.ErrorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}
			http.SetCookie(w, app.Auth.GetRefreshCookie(tokenPairs.RefreshToken))
			w.Header().Set("Content-Type", "AuthServerApp/json")
			json.NewEncoder(w).Encode(tokenPairs)
			return
		}
	}
	utils.JSONResponse{}.ErrorJSON(w, errors.New("no more cookies"), http.StatusUnauthorized)
}

// ValidateSession checks the validity of the JWT token provided in the Authorization header.
// It expects the token to be in the "Bearer <token>" format.
// If the token is valid and not expired, it responds with a success message.
// If the token is missing, invalid, or expired, it responds with an error message and appropriate HTTP status code.
func (app *AuthServerApp) ValidateSession(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("missing Authorization header"), http.StatusUnauthorized)
		return
	}

	// Expect header format: "Bearer <token>"
	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("invalid Authorization header format"), http.StatusUnauthorized)
		return
	}
	tokenString := authHeader[len(prefix):]

	// Parse and validate the token
	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(app.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("invalid or expired token"), http.StatusUnauthorized)
		return
	}

	// Check expiration (optional, as jwt.ParseWithClaims already does this)
	if claims.ExpiresAt < time.Now().Unix() {
		utils.JSONResponse{}.ErrorJSON(w, errors.New("token expired"), http.StatusUnauthorized)
		return
	}

	// If valid, return success
	resp := utils.JSONResponse{
		Error:   false,
		Message: "session is valid",
	}
	utils.JSONResponse{}.WriteJSON(w, http.StatusOK, resp)

}

// Logout invalidates the user's session by setting an expired refresh token cookie.
// This effectively logs the user out by removing the ability to refresh the JWT token.
// It responds with an HTTP 202 Accepted status to indicate the logout request was processed.
func (app *AuthServerApp) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.Auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}
