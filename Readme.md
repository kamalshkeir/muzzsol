# muzzsol


## Setup

After cloning the repos

#### One line command mysql db
```sh
docker run --name mysql -p 3306:3306 -e "MYSQL_DATABASE=muzzdb" -e "MYSQL_USER=user" -e "MYSQL_PASSWORD=strongPass" -e "MYSQL_ROOT_PASSWORD=rootPass" -d  mysql:latest
```
 or create a local db
- rename `.env.example` to `.env`
- update `.env` DB infos , these infos will be loaded in struct `settings.Config` using [kenv](https://github.com/kamalshkeir/kenv)
- run `go run main.go`

The first time you run the project, it will auto migrate models using struct, i used my ORM [korm](https://github.com/kamalshkeir/korm) . 
So the same code will work with other sql dialects

#### CreateUser:
- `GET /user/create`

you will notice password hash the name as a password, so you can copy the name of any user and tested.

#### Login:
- `POST /login`
- email : name@email.com
- password: name 
- session encrypted using AES, i prefered over JWT and once implemented it's very straigtforward 

this will return a token , and add it to cookies , you are authenticated, you are ready for `/profiles`
you can change cookie session name using middlewares.CookieSessionName
#### Profiles
- `GET /profiles?age=x&gender=x`
this endpoint take only queryParams and get userId from request context passed through the AuthMiddleware

profiles are sorted by distance from the user_id
it should look like this:
```json
{
  "results": [
    {
      "id": 30,
      "name": "qxsjs",
      "gender": "male",
      "age": 52,
      "distanceFromMe": 1.7557933930191036
    },
    {
      "id": 5,
      "name": "vjhyf",
      "gender": "male",
      "age": 53,
      "distanceFromMe": 10.475076919529819
    },
    {
      "id": 20,
      "name": "dmyvl",
      "gender": "male",
      "age": 52,
      "distanceFromMe": 15.032404929858311
    },
    {
      "id": 6,
      "name": "mmzih",
      "gender": "male",
      "age": 56,
      "distanceFromMe": 19.181882524682916
    },
    {
      "id": 12,
      "name": "qfeml",
      "gender": "male",
      "age": 51,
      "distanceFromMe": 25.979435191333025
    },
    {
      "id": 11,
      "name": "pekad",
      "gender": "male",
      "age": 55,
      "distanceFromMe": 26.206583779720727
    }
  ]
}
```
#### Swipe
- finaly `POST /swipe` expect in body:
- ``profileId`` : profile.Id
- ``preference`` : 'yes' or 'no'
- return matched only if both swiper yes

after swipe it should delete the profile and not see it again, you can check ``/profiles`` to verify

#### To Not Miss
file `handlers/util.go` , these 2 functions took me some time, and remind me when i was studying mechanics :
- GenerateRandomLocation
- DistanceBetweenLocations


