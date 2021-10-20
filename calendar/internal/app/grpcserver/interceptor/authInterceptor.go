package interceptor

import (
	"calendar/internal/app"
	"calendar/internal/auth"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

type AuthInterceptor struct {
	server *app.IServer
}

func NewAuthInterceptor(Server *app.IServer) *AuthInterceptor {
	return &AuthInterceptor{
		server: Server,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary i: ", info.FullMethod)

		ctx, err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func newWrappedStream(stream grpc.ServerStream, ctx context.Context) *wrappedStream {
	s := &wrappedStream{
		ctx: ctx,
	}
	s.ServerStream = stream
	return s
}

func (s *wrappedStream) Context() context.Context {
	return s.ctx
}

func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream i: ", info.FullMethod)

		ctx, err := i.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		stream = newWrappedStream(stream, ctx)

		return handler(srv, stream)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	if i.accessibleRoute(method) {
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, app.ErrEmptyAuthToken)
	}

	accessToken := values[0]
	claims, err := i.server.JWTWrapper.ValidateToken(accessToken)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, app.ErrInvalidAuthToken, err)
	}

	userSession, err := i.server.GetUserSession(claims.ID)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, app.ErrSessionNotFound)
	}

	if _, ok := userSession[accessToken]; !ok {
		return ctx, status.Errorf(codes.Unauthenticated, app.ErrInvalidAuthToken)
	}

	i.server.Logger.Infof("User ID #%v auth by token", claims.ID)
	ctxUserID := context.WithValue(ctx, auth.KeyUserID, claims.ID)

	return ctxUserID, nil
}

func (i AuthInterceptor) accessibleRoute(method string) bool {
	accessibleRoutes := map[string]bool{
		"/auth.AuthService/Login": true,
	}

	return accessibleRoutes[method]
}
