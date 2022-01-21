package main

import (
	"strconv"
	"time"
)

/*
Redis has a feature that automatically deletes data whose expiration time has reached.
Redis can also handle a lot of writes and can scale horizontally.
Since Redis is a key-value storage, its keys need to be unique, to achieve this, we will use uuid as the key and use the user id as the value.
*/

// SetTokenMetadataToRedis :
/* We passed in the TokenDetails which have information about the expiration time of the JWTs and the uuids used when creating the JWTs.
If the expiration time is reached for either the refresh token or the access token, the JWT is automatically deleted from Redis.
*/
func SetTokenMetadataToRedis(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

// SearchMetadataInRedis : Search JWT Metadata Still exist in redis store
func SearchMetadataInRedis(authD *TokenMetadata) (uint64, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

// DeleteMetadataFromRedis : Delete Metadata from Redis
func DeleteMetadataFromRedis(givenUuid string) (int64, error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
