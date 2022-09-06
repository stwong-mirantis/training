package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestUserResource_DoesAuthTokenExist(t *testing.T) {
	ur := NewUserResource()
	assert.Equal(t, ur.DoesAuthTokenExist("token-hi"), false)
	username := "julip"
	user, err := ur.AddUser(username)
	require.NoError(t, err)
	assert.Equal(t, ur.DoesAuthTokenExist(user.UUID), true)
}

func TestUserResource_SetUserOnlineStatus(t *testing.T) {
	ur := NewUserResource()
	assert.Equal(t, ur.DoesAuthTokenExist("token-hi"), false)
	username := "julip"
	user, err := ur.AddUser(username)
	require.NoError(t, err)
	onlineStatus := false
	ur.SetUserOnlineStatusNoLock(&onlineStatus, user.UUID)
	user2 := ur.GetUserWithToken(user.UUID)
	assert.Equal(t, user2.OnlineStatus, &onlineStatus)
}

func TestUserResource_RemoveUser(t *testing.T) {
	ur := NewUserResource()
	username := "julip"
	expectedUser, err := ur.AddUser(username)
	require.NoError(t, err)
	onlineUsersBefore := ur.GetAllOnlineUsers()
	assert.Equal(t, ur.DoesAuthTokenExist(expectedUser.UUID), true, "user exists should be true")
	require.Len(t, onlineUsersBefore, 1)
	ur.RemoveUser(expectedUser.UUID)
	onlineUsersAfter := ur.GetAllOnlineUsers()
	assert.Equal(t, ur.DoesAuthTokenExist(expectedUser.UUID), false, "user exists should be false")
	require.Len(t, onlineUsersAfter, 0)
}

func TestUserResource_UpdateUserLastSeenTime(t *testing.T) {
	ur := NewUserResource()
	username := "julip"
	user, err := ur.AddUser(username)
	initialTime := user.LastSeenTime
	require.NoError(t, err)
	ur.UpdateUserLastSeenTime(user.UUID)
	user2, err := ur.GetUserWithUsername(username)
	require.NoError(t, err)
	assert.NotEqual(t, initialTime, user2.LastSeenTime)
}

func TestUserResource_GetUserWithUsername(t *testing.T) {
	ur := NewUserResource()
	username := "julip"
	user, err := ur.AddUser(username)
	require.NoError(t, err)
	user2, err := ur.GetUserWithUsername(user.Username)
	require.NoError(t, err)
	assert.Equal(t, user2, user)
}

func TestUserResource_GetUserWithToken(t *testing.T) {
	ur := NewUserResource()
	username := "julip"
	user, err := ur.AddUser(username)
	require.NoError(t, err)
	assert.Equal(t, ur.GetUserWithToken(user.UUID), user)
}

func TestUserResource_AddUser(t *testing.T) {
	ur := NewUserResource()
	onlineUsersBefore := ur.GetAllOnlineUsers()
	require.Len(t, onlineUsersBefore, 0)
	username := "julip"
	user, err := ur.AddUser(username)
	require.NoError(t, err)
	onlineUsersAfter := ur.GetAllOnlineUsers()
	require.Len(t, onlineUsersAfter, 1)
	assert.Equal(t, ur.DoesAuthTokenExist(user.UUID), true, "auth token should exist")
}

func TestUserResource_GetAllOnlineUsers(t *testing.T) {
	ur := NewUserResource()
	var expectedUsers []User
	for i := 0; i < 3; i++ {
		u, err := ur.AddUser(strconv.Itoa(i))
		require.NoError(t, err)
		expectedUsers = append(expectedUsers, u)
	}
	onlineUsers := ur.GetAllOnlineUsers()

	require.ElementsMatch(t, expectedUsers, onlineUsers)
}
