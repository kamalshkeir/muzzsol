package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/kamalshkeir/korm"
	"github.com/kamalshkeir/muzzsol/middlewares"
	"github.com/kamalshkeir/muzzsol/services"
	"github.com/kamalshkeir/muzzsol/types"
	"github.com/labstack/echo/v4"
)

var (
	// ErrInvalidInputs is error when user send wrong inputs
	ErrInvalidInputs = errors.New("invalid inputs")
	// ErrUserDoesNotExist is error when user not found in db
	ErrUserNotFound = errors.New("user not found")
	// ErrBadPassword is error when pass is wrong
	ErrBadPassword = errors.New("bad password")
)


// UserRandom handler /user/create , generate and return random user and location
func UserRandom(c echo.Context) error {
	// generate random user
	u := generateRandomUser()
	l := GenerateRandomLocation(nil,5)
	// insert user 
	_, err := korm.Model[types.User]().Insert(u)
	if err != nil {
		c.JSON(400,map[string]any{
			"error":err.Error(),
		})
		return nil
	}
	// get created user
	createdUser, err := korm.Model[types.User]().Where("name = ?", u.Name).One()
	if err != nil {
		c.JSON(400,map[string]any{
			"error":err.Error(),
		})
		return nil
	}
	// create location
	l.UserId=createdUser.Id
	_,err = korm.Model[types.Location]().Insert(&l)
	if err != nil {
		c.JSON(400,map[string]any{
			"error":err.Error(),
		})
		return nil
	}
	// return created user
	c.JSON(200, map[string]any{
		"result": createdUser,
	})
	return nil
}

// GetUserProfiles get profiles of specific user, user is accessible through the request ctx or getSessionUserId(c)
func GetUserProfiles(c echo.Context) error {
	// getting profile query param
	profileParams := types.ProfileParams{}
	if err := c.Bind(&profileParams); err != nil {
		return err
	}
	// getting user from session
	idUser := getSessionUserId(c)
	if idUser == -1 {
		c.JSON(200,map[string]any{
			"error":"Unauthorised access",
		})
		return nil
	} 
	// get user
	user, err := korm.Model[types.User]().Where("id = ?", idUser).One()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	// here i will get all users with opposite gender, and age higher or less than 10 years with min of age 18
	whereQuery := "id != ? AND gender != ? AND age > ? AND age < ? AND id NOT IN (select profile_id from swipes where user_id=?)"
	values := []any{user.Id, user.Gender, math.Max(18,float64(user.Age-10)) , user.Age+10, user.Id}
	if profileParams.Age != "" {
		whereQuery += " and age = ?"
		age,err := strconv.Atoi(profileParams.Age)
		if err != nil {
			return err
		}
		values = append(values, age)
	}

	// here filter by gender doesn't make sense because i always take the opposite gender, but still work, you will get no data if same gender
	if profileParams.Gender != "" {
		whereQuery += " and gender = ?"
		values = append(values, profileParams.Gender)
	}
	// get profiles
	profiles, err := korm.Model[types.User]().Select("id", "name", "gender", "age").Where(whereQuery,values...).All()
	if err != nil {
		c.JSON(200,map[string]any{
			"error":err.Error(),
		})
		return nil
	}

	// get user location
	userLocation,err := korm.Model[types.Location]().Where("user_id = ?",user.Id).One()
	if err != nil && err.Error() != "no data found" {
		c.Logger().Error(err)
	} else {
		calculateAndSortByDistance(userLocation,profiles)
	}
	// send results
	c.JSON(200, map[string]any{
		"results": profiles,
	})
	return nil
}


// UserSwipe handle a user swipe, and update list of potential matches
func UserSwipe(c echo.Context) error {
	// get user from session
	idUser := getSessionUserId(c)
	if idUser == -1 {
		c.JSON(200,map[string]any{
			"error":"Unauthorised access",
		})
		return nil
	}
	swipeParams := types.SwipeParams{}
	if err := c.Bind(&swipeParams); err != nil {
		return err
	}
	// check if user exist
	user, err := korm.Model[types.User]().Where("id = ?", idUser).One()
	if err != nil {
		c.JSON(http.StatusBadRequest,map[string]any{
			"error":ErrUserNotFound.Error(),
		})
		return nil
	}
	// get profile id from body
	preference := swipeParams.Preference
	if preference != "yes" && preference != "no" {
		c.JSON(http.StatusBadRequest,map[string]any{
			"error":"expected preference to be 'yes' or 'no'",
		})
		return nil
	}
	profileId, err := strconv.Atoi(swipeParams.ProfileId)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	// store the swipe, check if exist before
	var profileSwipe types.Swipe
	profileSwipe, err = korm.Model[types.Swipe]().Where("user_id = ? and profile_id = ?", user.Id, profileId).One()
	if err == nil {
		c.JSON(http.StatusBadRequest,map[string]any{
			"error":"already swiped profil",
		})
		return nil
	}
	swiper := types.Swipe{
		UserId:     uint(user.Id),
		ProfileId:  uint(profileId),
		Preference: swipeParams.Preference,
	}
	_, err = korm.Model[types.Swipe]().Insert(&swiper)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	// check if there is a match
	profileSwipe, err = korm.Model[types.Swipe]().Where("user_id = ? and profile_id = ?", profileId, user.Id).One()
	// if the profile did not swipe or swipe left no match yet
	if err != nil || profileSwipe.Preference == "no" {
		c.JSON(200, map[string]any{
			"matched": false,
		})
		return nil
	}
	// we have a match, we return matchID the id of matched profile
	c.JSON(200, map[string]any{
		"results": map[string]any{
			"matched": true,
			"matchID": profileSwipe.Id,
		},
	})
	return nil
}

// Login authenticate a user, set a cookie named session
func Login(c echo.Context) error {
	loginParams := types.LoginParams{}
	if err := c.Bind(&loginParams); err != nil {
		return err
	}
	// validate email
	_, err := mail.ParseAddress(loginParams.Email)
	if err != nil || loginParams.Password == "" {
		c.JSON(http.StatusBadRequest,map[string]any{
			"error":ErrInvalidInputs.Error(),
		})
		return nil
	}
	// get user from email
	user,err := korm.Model[types.User]().Where("email = ?", loginParams.Email).One()
	if err != nil  {
		c.JSON(http.StatusBadRequest,map[string]any{
			"error":ErrUserNotFound.Error(),
		})
		return nil
	}
	// hash input password and compare it to db
	valid,err := services.ComparePasswordToHash(loginParams.Password,user.Password)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest,map[string]any{
			"error":ErrBadPassword.Error(),
		})
		return nil
	}
	cookieData,err := json.Marshal(map[string]any{
		"id":user.Id,
	})
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	// puting the user id in cookies can be sensitive, that's why i choose to use AES encryption, because it add extra security by not allowing access to the content of cookies , https://medium.com/@arpitkh96/is-jwt-authentication-enough-751699b0f58d

	// encrypt and set session cookie
	enc,err := services.Encrypt(string(cookieData))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	// store cookie
	c.SetCookie(&http.Cookie{
		Name: middlewares.CookieSessionName,
		Value: enc,
		Path: "/",
		Expires: time.Now().Add(24*time.Hour*7),
		HttpOnly: true,
	})
	// return token
	c.JSON(200,map[string]any{
		"token":enc,
	})
	return nil
}




