## bootstrap
* provide config to enable initial realm provisioning with admin user

## pagination
* cursor pagination

## subprojects
* tribe-api (plugins api / database api / flows )
* tribe-ui (user interface)

## database design

realm
* realm_id
* name (optional) 

user
* id unique ID  (unique - globally)

* realm_id (default main - configurable name)
* username (unique - in realm)

[basic profile](https://openid.net/specs/openid-connect-basic-1_0.html)

* display_name (name)
* given_name ()
* family_name
* middle_name
* nickname
* preferred_username 
* profile
* picture
* website
* email
* email_verified

## address

## TODO
- [ ] POST /signup 
    - Register a new user with an email and password
- [ ] POST /invite
    - Invites a new user with an email.
- [ ] POST /verify
    - Verify a registration or a password recovery.
- [ ] POST /recover  
    - Password recovery. Will deliver a password recovery mail to the user based on email address.  
- [ ] POST /token
    - This is an OAuth2 endpoint that currently implements the password, refresh_token, and authorization_code grant types
- [ ] GET /user
    - Get the JSON object for the logged in user (requires authentication)
- [ ] PUT /user
    - Update a user (Requires authentication). 
      Apart from changing email/password, this method can be used to set custom user data.
- [ ] POST /logout
    - Logout a user (Requires authentication).
    
