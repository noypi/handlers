package handlers

import (
	"fmt"
	"net/http"

	"errors"

	"github.com/noypi/logfn"
	"github.com/noypi/router"

	"github.com/huandu/facebook"
)

type _FbUserAuthenticatedName int

const FbUserAuthenticatedName _FbUserAuthenticatedName = 0

type FacebookMe struct {
	Name, Id string
}

func ValidateFbUser(userField, accessTokenField string, urlFetchClientKeyName interface{}) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			err = recover2Err(recover(), err)
			if nil != err {
				//w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
		}()

		ctx := router.ContextW(w)

		INFO := ctx.MustGet("INFO").(logfn.LogFunc)
		ERR := ctx.MustGet("ERR").(logfn.LogFunc)

		INFO("+ValidateFbUser()")
		defer INFO("-ValidateFbUser()")

		fbme, bValid := isFbUserReal(w, r, userField, accessTokenField, urlFetchClientKeyName)
		if !bValid {
			w.WriteHeader(http.StatusUnauthorized)
			ERR.PrintStackTrace(5)
			err = fmt.Errorf("ValidateFbUser: invalid fbuser")
			return
		}

		ctx.Set(FbUserAuthenticatedName, fbme)
	}

}

func isFbUserReal(w http.ResponseWriter, r *http.Request, userField, accessTokenField string, urlFetchClientKeyName interface{}) (fbme *FacebookMe, bValid bool) {
	ctx := router.ContextW(w)

	INFO := ctx.MustGet("INFO").(logfn.LogFunc)
	ERR := ctx.MustGet("ERR").(logfn.LogFunc)

	client := ctx.MustGet(urlFetchClientKeyName).(*http.Client)
	if nil == client {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ERR.PrintStackTrace(5)
		ERR("ValidateFbUser: net-client is not found.")
		return
	}

	r.ParseForm()
	fbId := r.PostFormValue(userField)
	fbAccessToken := r.PostFormValue(accessTokenField)

	facebook.SetHttpClient(client)
	//verify
	res, err := facebook.Get("/me", facebook.Params{
		"access_token": fbAccessToken,
	})
	if nil != err {
		// wrong
		ERR.PrintStackTrace(5)
		ERR("fb.Get err=%v", err)
		return
	}

	fbme = new(FacebookMe)
	if err = res.Decode(fbme); nil != err {
		ERR.PrintStackTrace(5)
		ERR("res.Decode err=%v", err)
		return
	}

	INFO("fbme.Name=%s, fbme.Id=%s, curr fbid=%s", fbme.Name, fbme.Id, fbId)

	bValid = (fbme.Id == fbId)
	return
}

func recover2Err(a interface{}, def error) error {
	switch v := a.(type) {
	case nil:
		return nil
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return def
	}
}
