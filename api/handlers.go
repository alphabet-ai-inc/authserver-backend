package api

import (
	"backend/auth"
	"backend/internal/dbrepo"
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

type Autserverapp struct {
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

func (app *Autserverapp) Home(w http.ResponseWriter, r *http.Request) {

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

func (app *Autserverapp) Apps(w http.ResponseWriter, r *http.Request) {

	apps, err := app.DB.AllApps()
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	_ = utils.JSONResponse.WriteJSON(utils.JSONResponse{}, w, http.StatusOK, apps)
	// if err != nil {
	// 	http.Error(w, "Unable to fetch apps", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "Autserverapp/json")
	json.NewEncoder(w).Encode(apps)

}

func (app *Autserverapp) AppsCatalogue(w http.ResponseWriter, r *http.Request) {

	apps, err := app.DB.AllApps()
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}
	_ = utils.JSONResponse{}.WriteJSON(w, http.StatusOK, apps)

}

func (app *Autserverapp) GetApp(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	appID, err := strconv.Atoi(id)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}
	thisapp, err := app.DB.ThisApp(appID)
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	_ = utils.JSONResponse{}.WriteJSON(w, http.StatusOK, thisapp)
}

func (app *Autserverapp) ThisApp(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	// id := chi.URLParam(r, "id")

	appID, err := strconv.Atoi(id)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	thisapp, err := app.DB.ThisApp(appID)
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	_ = utils.JSONResponse{}.WriteJSON(w, http.StatusOK, thisapp)
}

func (app *Autserverapp) ThisAppForEdit(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	appID, err := strconv.Atoi(id)
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	thisapp, err := app.DB.ThisApp(appID)
	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	_ = utils.JSONResponse{}.WriteJSON(w, http.StatusOK, thisapp)
}

func (app *Autserverapp) InsertApp(w http.ResponseWriter, r *http.Request) {
	var thisapp models.ThisApp

	err := utils.JSONResponse{}.ReadJSON(w, r, &thisapp)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}
	var newapp models.NewApp
	newapp.Name = mux.Vars(r)["name"]
	newapp.Release = thisapp.Release
	newapp.Path = thisapp.Path
	newapp.Init = thisapp.Init
	newapp.Web = thisapp.Web
	newapp.Title = thisapp.Title
	newapp.Created = time.Now().Unix()
	newapp.Updated = time.Now().Unix()

	newID, err := app.DB.InsertApp(newapp)
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
func (app *Autserverapp) UpdateApp(w http.ResponseWriter, r *http.Request) {
	var payload models.ThisApp

	err := utils.JSONResponse{}.ReadJSON(w, r, &payload)

	if err != nil {
		utils.JSONResponse{}.ErrorJSON(w, err)
		return
	}

	thisapp, err := app.DB.ThisApp(payload.ID)

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
	thisapp.Created = time.Now().Unix()
	thisapp.Updated = time.Now().Unix()

	err = app.DB.UpdateApp(*thisapp)

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

func (app *Autserverapp) DeleteApp(w http.ResponseWriter, r *http.Request) {

	// id := chi.URLParam(r, "id")
	id := mux.Vars(r)["id"]

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
func (app *Autserverapp) Authenticate(w http.ResponseWriter, r *http.Request) {
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

	refreshCookie := app.Auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	utils.JSONResponse{}.WriteJSON(w, http.StatusAccepted, tokens)
	w.Header().Set("Content-Type", "Autserverapp/json")
	json.NewEncoder(w).Encode(tokens)

}

func (app *Autserverapp) RefreshToken(w http.ResponseWriter, r *http.Request) {
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
			w.Header().Set("Content-Type", "Autserverapp/json")
			json.NewEncoder(w).Encode(tokenPairs)
			return
		}
	}
	utils.JSONResponse{}.ErrorJSON(w, errors.New("no more cookies"), http.StatusUnauthorized)
}

func (app *Autserverapp) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.Auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}
