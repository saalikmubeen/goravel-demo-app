package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/saalikmubeen/goravel"
	"github.com/saalikmubeen/goravel-demo-app/models"
	"github.com/saalikmubeen/goravel/mailer"
	"github.com/saalikmubeen/goravel/urlsigner"
)

// UserLogin displays the login page
func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// PostUserLogin attempts to log a user in
func (h *Handlers) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	matches, err := user.PasswordMatches(password)
	if err != nil {
		w.Write([]byte("Error validating password"))
		return
	}

	if !matches {
		w.Write([]byte("Invalid password!"))
		return
	}

	// did the user check remember me?
	if r.Form.Get("remember") == "remember" {

		rm := models.RememberMeToken{}
		sha, err := rm.InsertToken(user.ID)
		if err != nil {
			h.App.ErrorStatus(w, http.StatusBadRequest)
			return
		}

		// set a cookie
		expire := time.Now().Add(365 * 24 * 60 * 60 * time.Second) // 1 year
		cookie := http.Cookie{
			Name:     fmt.Sprintf("_%s_remember_me", h.App.AppName),
			Value:    fmt.Sprintf("%d|%s", user.ID, sha), // user_id|remember_me_token
			Path:     "/",
			Expires:  expire,
			HttpOnly: true,
			Domain:   h.App.Session.Cookie.Domain,
			MaxAge:   31536000, // 1 years
			Secure:   h.App.Session.Cookie.Secure,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		// save the remember_me_token in the user's session
		h.App.Session.Put(r.Context(), "remember_me_token", sha)
	}

	h.App.Session.Put(r.Context(), "userID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Logout logs the user out, removes any remember me cookie, and deletes
// remember token from the database, if it exists
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	// delete the remember token if it exists
	if h.App.Session.Exists(r.Context(), "remember_token") {
		rm := models.RememberMeToken{}

		// Get the remember_me_token from the user's session
		rememberMeToken := h.App.Session.GetString(r.Context(), "remember_me_token")
		_ = rm.Delete(rememberMeToken)
	}

	// delete cookie
	newCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_me_remember", h.App.AppName),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain:   h.App.Session.Cookie.Domain,
		MaxAge:   -1,
		Secure:   h.App.Session.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &newCookie)

	h.App.Session.RenewToken(r.Context())               // to prevent session fixation
	h.App.Session.Remove(r.Context(), "userID")         // remove the user's ID from the users's session stored in the session store
	h.App.Session.Remove(r.Context(), "remember_token") // remove the remember_me_token from the user's session stored in the session store
	h.App.Session.Destroy(r.Context())                  // remove or delete the whole session entry corresponding to the user from the session store
	h.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (h *Handlers) UserSignup(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "signup", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUserSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")

	// validate the form
	validator := h.App.Validator(r.Form)

	validator.IsValidEmail("email")
	validator.IsValidPassword("password")
	validator.HasMinLength("first_name", 4)
	validator.HasMinLength("last_name", 4)

	if !validator.IsValid() {
		vars := make(jet.VarMap)
		vars.Set("validator", validator)
		var user models.User
		user.FirstName = first_name
		user.LastName = last_name
		user.Email = email
		user.Password = password
		vars.Set("user", user)

		if err := h.App.Render.Page(w, r, "signup", vars, nil); err != nil {
			h.App.ErrorLog.Println(err)
			return
		}
		return
	}

	user := models.User{
		Email:     email,
		Password:  password,
		FirstName: first_name,
		LastName:  last_name,
	}

	_, err = user.Insert(user)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (h *Handlers) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "forgot-password", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error rendering: ", err)
		h.App.Error500(w, r)
	}
}

// PostForgot looks up a user by email, and if the user is found, generates
// and sends an email with a singed link to the reset password form
func (h *Handlers) PostForgotPassword(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorStatus(w, http.StatusBadRequest)
		return
	}

	// verify that supplied email exists
	var u *models.User
	email := r.Form.Get("email")
	u, err = u.GetByEmail(email)
	if err != nil {
		h.App.ErrorStatus(w, http.StatusBadRequest)
		return
	}

	// create a link to password reset form
	link := fmt.Sprintf("%s/users/reset-password?email=%s", h.App.Server.URL, email)

	// sign the link
	sign := urlsigner.Signer{
		Secret: []byte(h.App.EncryptionKey),
	}

	signedLink := sign.GenerateTokenFromString(link)
	h.App.InfoLog.Println("Signed link is", signedLink)

	// email the message
	var data struct {
		Link string
	}
	data.Link = signedLink

	msg := mailer.Message{
		To:       u.Email,
		Subject:  "Password reset",
		Template: "password-reset",
		Data:     data,
		From:     "salikmubien@gmail.com",
	}

	h.App.Mail.Jobs <- msg
	res := <-h.App.Mail.Results
	if res.Error != nil {
		fmt.Println("Error sending email: ", res.Error)
		h.App.ErrorStatus(w, http.StatusBadRequest)
		return
	}

	// redirect the user
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

// ResetPasswordForm validates a signed url, and displays the password reset form, if appropriate
func (h *Handlers) ResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	// get form values
	email := r.URL.Query().Get("email")
	theURL := r.RequestURI
	// fmt.Println(theURL)           // /users/reset-password?email=salikmubien@gmail.com&hash=.3BLce3.q04wGf1KS1X8JX3FIW4WTvlzKq3n3S-zAHBDr8uDOGg
	// fmt.Println(h.App.Server.URL) // http://localhost:4000
	testURL := fmt.Sprintf("%s%s", h.App.Server.URL, theURL)

	// validate the url
	signer := urlsigner.Signer{
		Secret: []byte(h.App.EncryptionKey),
	}

	valid := signer.VerifyToken(testURL)
	if !valid {
		h.App.ErrorLog.Print("Invalid url")
		h.App.ErrorUnauthorized(w, r)
		return
	}

	// make sure it's not expired
	// less than 60 minutes old
	expired := signer.Expired(testURL, 60)
	if expired {
		h.App.ErrorLog.Print("Link expired")
		h.App.ErrorUnauthorized(w, r)
		return
	}

	enc := goravel.Encryption{Key: []byte(h.App.EncryptionKey)}

	// display form
	encryptedEmail, _ := enc.Encrypt(email)

	vars := make(jet.VarMap)
	vars.Set("email", encryptedEmail)

	err := h.App.Render.Page(w, r, "reset-password", vars, nil)
	if err != nil {
		return
	}
}

// PostResetPassword resets the user's password
func (h *Handlers) PostResetPassword(w http.ResponseWriter, r *http.Request) {
	// parse the form
	err := r.ParseForm()
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	enc := goravel.Encryption{Key: []byte(h.App.EncryptionKey)}

	// get and decrypt the email
	email, err := enc.Decrypt(r.Form.Get("email"))
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	// get the user
	var u models.User
	user, err := u.GetByEmail(email)
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	// reset the password
	err = user.ResetPassword(user.ID, r.Form.Get("password"))
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	// redirect
	h.App.Session.Put(r.Context(), "flash", "Password reset. You can now log in.")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (h *Handlers) CurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	// get the user
	id := h.App.Session.GetInt(r.Context(), "userID")

	if id == 0 {
		h.App.ErrorStatus(w, http.StatusUnauthorized)
		return
	}

	var u models.User
	user, err := u.Get(id)
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("User", user)

	err = h.App.Render.Page(w, r, "user-profile", vars, nil)
	if err != nil {
		h.App.Error500(w, r)
		return
	}
}
