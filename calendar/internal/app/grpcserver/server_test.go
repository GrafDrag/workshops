package grpcserver_test

import (
	"calendar/internal/app/grpcserver"
	"calendar/internal/auth"
	"calendar/internal/config"
	"calendar/internal/model"
	sess "calendar/internal/session/inmemory"
	"calendar/internal/store/inmemory"
	"calendar/pb"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const bufSize = 1024 * 1024

var (
	store      = inmemory.New()
	jwtWrapper = &auth.JwtWrapper{
		SecretKey:       "secretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}
	session = sess.NewSession(config.SessionConfig{})
	lis     *bufconn.Listener
)

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpcserver.NewServer(store, jwtWrapper, session)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("server exited with error: %v", err)
		}
	}()
}

func TestUser_Login(t *testing.T) {
	t.Run("test user login", func(t *testing.T) {
		user := model.TestUser(t)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		client := pb.NewAuthServiceClient(conn)
		resp, err := client.Login(ctx, &pb.LoginRequest{
			Login:    user.Login,
			Password: user.Password,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})
}

func TestUser_Logout(t *testing.T) {
	t.Run("test user logout", func(t *testing.T) {
		user := model.TestUser(t)
		token := getJWTToken(t, user)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		client := pb.NewUserServiceClient(conn)
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx = metadata.NewOutgoingContext(ctx, md)

		resp, err := client.Logout(ctx, &pb.LogoutRequest{})
		assert.NoError(t, err)
		assert.Equal(t, resp.Status, pb.LogoutResponse_Successful)

		_, err = client.Logout(ctx, &pb.LogoutRequest{})
		assert.Error(t, err)
	})
}

func TestUser_Update(t *testing.T) {
	t.Run("update user timezone", func(t *testing.T) {
		timezone := "Europe/Kiev"
		user := model.TestUser(t)
		token := getJWTToken(t, user)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		client := pb.NewUserServiceClient(conn)
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx = metadata.NewOutgoingContext(ctx, md)

		resp, err := client.Update(ctx, &pb.UserUpdateRequest{
			Login:    user.Login,
			Timezone: timezone,
		})

		assert.NoError(t, err)
		assert.Equal(t, resp.Status, pb.UserUpdateResponse_Successful)

		got, _ := store.User().FindById(user.ID)
		assert.Equal(t, timezone, got.Timezone)
	})

	t.Run("update user login", func(t *testing.T) {
		newLogin := "testUser2"
		user := model.TestUser(t)
		token := getJWTToken(t, user)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		client := pb.NewUserServiceClient(conn)
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx = metadata.NewOutgoingContext(ctx, md)

		resp, err := client.Update(ctx, &pb.UserUpdateRequest{
			Login:    newLogin,
			Timezone: user.Timezone,
		})

		assert.NoError(t, err)
		assert.Equal(t, resp.Status, pb.UserUpdateResponse_Successful)

		got, _ := store.User().FindById(user.ID)
		assert.Equal(t, newLogin, got.Login)
	})
}

func TestServer_ListEvents(t *testing.T) {
	user := model.TestUser(t)
	token := getJWTToken(t, user)
	event := model.TestEvent(t)
	event.UserID = user.ID
	event.Timezone = "Europe/Kiev"
	addEventToStore(t, event)
	var cases = []struct {
		name    string
		search  func() *pb.ListRequest
		isEmpty bool
	}{
		{
			name: "search by title",
			search: func() *pb.ListRequest {
				return &pb.ListRequest{
					Title: event.Title,
				}
			},
			isEmpty: false,
		},
		{
			name: "search by timezone",
			search: func() *pb.ListRequest {
				return &pb.ListRequest{
					Timezone: event.Timezone,
				}
			},
			isEmpty: false,
		},
		{
			name: "search by dateFrom",
			search: func() *pb.ListRequest {
				return &pb.ListRequest{
					DateFrom: time.Now().Format("2006-01-02"),
				}
			},
			isEmpty: false,
		},
		{
			name: "search by dateTo",
			search: func() *pb.ListRequest {
				return &pb.ListRequest{
					DateTo: time.Now().Format("2006-01-02"),
				}
			},
			isEmpty: true,
		},
		{
			name: "search by dateFrom with timeFrom",
			search: func() *pb.ListRequest {
				now := time.Now().Add(time.Minute * -10)
				return &pb.ListRequest{
					DateFrom: now.Format("2006-01-02"),
					TimeFrom: now.Format("15:04"),
				}
			},
			isEmpty: false,
		},
		{
			name: "search by dateTo with timeTo",
			search: func() *pb.ListRequest {
				now := time.Now().Add(time.Minute * 20)
				return &pb.ListRequest{
					DateTo: now.Format("2006-01-02"),
					TimeTo: now.Format("15:04"),
				}
			},
			isEmpty: false,
		},
	}
	for _, s := range cases {
		t.Run(s.name, func(t *testing.T) {
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer func() {
				err := conn.Close()
				if err != nil {
					return
				}
			}()

			client := pb.NewEventServiceClient(conn)
			md := metadata.New(map[string]string{
				"authorization": token,
			})
			ctx = metadata.NewOutgoingContext(ctx, md)

			resp, err := client.List(ctx, s.search())

			assert.NoError(t, err)
			assert.Equal(t, s.isEmpty, len(resp.Event) == 0)
		})
	}
}

func TestServer_CreateEvent(t *testing.T) {
	t.Run("test event create", func(t *testing.T) {
		event := model.TestEvent(t)
		token := getJWTToken(t, model.TestUser(t))

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		client := pb.NewEventServiceClient(conn)
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx = metadata.NewOutgoingContext(ctx, md)

		resp, err := client.Create(ctx, &pb.CreateRequest{
			Title:       event.Title,
			Description: event.Description,
			Time:        event.Time,
			Timezone:    event.Timezone,
			Duration:    event.Duration,
			Notes:       event.Notes,
		})

		assert.NoError(t, err)
		assert.Equal(t, resp.Status, pb.CreateResponse_Successful)
	})
}

func TestServer_GetEventsById(t *testing.T) {
	t.Run("test get event by id", func(t *testing.T) {
		event := model.TestEvent(t)
		user := model.TestUser(t)
		token := getJWTToken(t, user)
		event.UserID = user.ID

		addEventToStore(t, event)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		client := pb.NewEventServiceClient(conn)
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx = metadata.NewOutgoingContext(ctx, md)

		resp, err := client.GetById(ctx, &pb.GetRequest{
			Id: int32(event.ID),
		})

		assert.NoError(t, err)
		assert.Equal(t, event.ID, int(resp.Event.Id))
	})
}

func TestServer_HandleUpdateEvent(t *testing.T) {
	cases := []struct {
		name  string
		field string
		value string
	}{
		{
			name:  "update title",
			field: "Title",
			value: "New title",
		},
		{
			name:  "update description",
			field: "Description",
			value: "New description",
		},
		{
			name:  "update time",
			field: "Time",
			value: time.Now().Add(50 * time.Minute).Format(model.EventDateLayout),
		},
		{
			name:  "update timezone",
			field: "Timezone",
			value: "Europe/Kiev",
		},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			event := model.TestEvent(t)
			user := model.TestUser(t)
			token := getJWTToken(t, user)
			event.UserID = user.ID
			addEventToStore(t, event)

			f := reflect.ValueOf(event).Elem().FieldByName(item.field)
			switch f.Kind() {
			case reflect.Int:
				val, err := strconv.Atoi(item.value)
				if err != nil {
					t.Fatal("failed convert str to int64")
				}
				f.SetInt(int64(val))
			case reflect.String:
				f.SetString(item.value)
			}

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer func() {
				err := conn.Close()
				if err != nil {
					return
				}
			}()

			client := pb.NewEventServiceClient(conn)
			md := metadata.New(map[string]string{
				"authorization": token,
			})
			ctx = metadata.NewOutgoingContext(ctx, md)

			resp, err := client.Update(ctx, &pb.UpdateRequest{
				Id:          int32(event.ID),
				Title:       event.Title,
				Description: event.Description,
				Time:        event.Time,
				Timezone:    event.Timezone,
				Duration:    event.Duration,
				Notes:       event.Notes,
			})

			assert.NoError(t, err)
			assert.Equal(t, resp.Status, pb.UpdateResponse_Successful)

			got, _ := store.Event().FindById(event.ID)

			assert.True(t, equals(got, event))
		})
	}
}

func TestServer_DeleteEvent(t *testing.T) {
	t.Run("delete event", func(t *testing.T) {
		event := model.TestEvent(t)
		user := model.TestUser(t)
		token := getJWTToken(t, user)
		event.UserID = user.ID
		addEventToStore(t, event)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		client := pb.NewEventServiceClient(conn)
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx = metadata.NewOutgoingContext(ctx, md)

		resp, err := client.Delete(ctx, &pb.DeleteRequest{
			Id: int32(event.ID),
		})

		assert.NoError(t, err)
		assert.Equal(t, resp.Status, pb.DeleteResponse_Successful)

		_, err = store.Event().FindById(event.ID)
		assert.Error(t, err)
	})
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func getJWTToken(t *testing.T, user *model.User) string {
	if err := store.User().Create(user); err != nil {
		t.Fatal("could not create user")
	}
	token, _ := jwtWrapper.GenerateToken(user)
	jb, err := json.Marshal(map[string]bool{
		token: true,
	})

	if err != nil {
		t.Fatal("could not marshal user map session")
	}

	if err := session.Set(strconv.Itoa(user.ID), string(jb)); err != nil {
		t.Fatal("could not save user session")
	}

	return token
}

func addEventToStore(t *testing.T, event *model.Event) {
	if err := store.Event().Create(event); err != nil {
		t.Fatal("could not create event")
	}
}

func equals(got, wont *model.Event) bool {
	for i, note := range got.Notes {
		if note != wont.Notes[i] {
			return false
		}
	}
	return got.Title == wont.Title &&
		got.Description == wont.Description &&
		got.Time == wont.Time &&
		got.Timezone == wont.Timezone &&
		got.Duration == wont.Duration &&
		len(got.Notes) == len(wont.Notes)
}
