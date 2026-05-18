package integration_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	protov1 "github.com/bhcoder23/gin-clean-template/docs/proto/v1"
	natsClient "github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/client"
	rmqClient "github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type notificationResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	TaskID    string `json:"task_id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
	ReadAt    string `json:"read_at"`
}

func TestHTTPNotificationsV1(t *testing.T) {
	token := registerAndLogin(t)
	httpCreateTask(t, token, "notify task", "creates notification")

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doAuthenticatedRequest(ctx, http.MethodGet, basePathV1+"/notifications?unread_only=true&limit=10&offset=0", http.NoBody, token)
	if err != nil {
		t.Fatalf("List notifications: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	type listResponse struct {
		Notifications []notificationResponse `json:"notifications"`
		Total         int                    `json:"total"`
	}

	listed := parseJSON[listResponse](t, resp)
	if listed.Total < 1 || len(listed.Notifications) == 0 {
		t.Fatalf("expected at least one notification, got total=%d", listed.Total)
	}

	notificationID := listed.Notifications[0].ID

	ctx2, cancel2 := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel2()

	markResp, err := doAuthenticatedRequest(ctx2, http.MethodPatch, basePathV1+"/notifications/"+notificationID+"/read", http.NoBody, token)
	if err != nil {
		t.Fatalf("Mark notification read: %v", err)
	}
	defer markResp.Body.Close()

	if markResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", markResp.StatusCode)
	}

	updated := parseJSON[notificationResponse](t, markResp)
	if !updated.Read {
		t.Fatal("expected notification to be marked as read")
	}
}

func TestGRPCNotificationsV1(t *testing.T) {
	token := registerAndLoginGRPC(t)

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	defer func() {
		if cerr := grpcConn.Close(); cerr != nil {
			t.Fatalf("grpcConn.Close: %v", cerr)
		}
	}()

	taskClient := protov1.NewTaskServiceClient(grpcConn)
	_, err = taskClient.CreateTask(grpcAuthCtx(t, token), &protov1.CreateTaskRequest{
		Title:       "grpc notify task",
		Description: "create notification",
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	notificationClient := protov1.NewNotificationServiceClient(grpcConn)
	listed, err := notificationClient.ListNotifications(grpcAuthCtx(t, token), &protov1.ListNotificationsRequest{
		UnreadOnly: true,
		Limit:      10,
		Offset:     0,
	})
	if err != nil {
		t.Fatalf("ListNotifications: %v", err)
	}

	if listed.GetTotal() < 1 || len(listed.GetNotifications()) == 0 {
		t.Fatalf("expected notifications, got total=%d", listed.GetTotal())
	}

	updated, err := notificationClient.MarkNotificationRead(grpcAuthCtx(t, token), &protov1.MarkNotificationReadRequest{
		Id: listed.GetNotifications()[0].GetId(),
	})
	if err != nil {
		t.Fatalf("MarkNotificationRead: %v", err)
	}

	if !updated.GetRead() {
		t.Fatal("expected notification to be marked as read")
	}
}

func TestRMQNotificationsV1(t *testing.T) {
	token := registerAndLoginRMQ(t)

	client, err := rmqClient.New(rmqURL, rpcServerExchange, rpcClientExchange)
	if err != nil {
		t.Fatalf("rmqClient.New: %v", err)
	}
	defer func() {
		if serr := client.Shutdown(); serr != nil {
			t.Fatalf("client.Shutdown: %v", serr)
		}
	}()

	createPayload := map[string]string{
		"title":       "rmq notify task",
		"description": "create notification",
	}

	var created taskResponse
	if err = client.RemoteCall("v1.task.create", authenticatedPayload(token, createPayload), &created); err != nil {
		t.Fatalf("task create: %v", err)
	}

	var listed struct {
		Notifications []notificationResponse `json:"notifications"`
		Total         int                    `json:"total"`
	}

	if err = client.RemoteCall("v1.notification.list", authenticatedPayload(token, map[string]any{
		"unread_only": true,
		"limit":       10,
		"offset":      0,
	}), &listed); err != nil {
		t.Fatalf("notification list: %v", err)
	}

	if listed.Total < 1 || len(listed.Notifications) == 0 {
		t.Fatalf("expected notifications after creating task %s", created.ID)
	}

	var updated notificationResponse
	if err = client.RemoteCall("v1.notification.markRead", authenticatedPayload(token, map[string]string{
		"id": listed.Notifications[0].ID,
	}), &updated); err != nil {
		t.Fatalf("notification markRead: %v", err)
	}

	if !updated.Read {
		t.Fatal("expected notification to be marked as read")
	}
}

func TestNATSNotificationsV1(t *testing.T) {
	token := registerAndLoginNATS(t)

	client, err := natsClient.New(natsURL, rpcServerExchange)
	if err != nil {
		t.Fatalf("natsClient.New: %v", err)
	}
	defer func() {
		if serr := client.Shutdown(); serr != nil {
			t.Fatalf("client.Shutdown: %v", serr)
		}
	}()

	var created taskResponse
	if err = client.RemoteCall("v1.task.create", authenticatedPayload(token, map[string]string{
		"title":       "nats notify task",
		"description": "create notification",
	}), &created); err != nil {
		t.Fatalf("task create: %v", err)
	}

	var listed struct {
		Notifications []notificationResponse `json:"notifications"`
		Total         int                    `json:"total"`
	}

	if err = client.RemoteCall("v1.notification.list", authenticatedPayload(token, map[string]any{
		"unread_only": true,
		"limit":       10,
		"offset":      0,
	}), &listed); err != nil {
		t.Fatalf("notification list: %v", err)
	}

	if listed.Total < 1 || len(listed.Notifications) == 0 {
		t.Fatalf("expected notifications after creating task %s", created.ID)
	}

	var updated notificationResponse
	if err = client.RemoteCall("v1.notification.markRead", authenticatedPayload(token, map[string]string{
		"id": listed.Notifications[0].ID,
	}), &updated); err != nil {
		t.Fatalf("notification markRead: %v", err)
	}

	if !updated.Read {
		t.Fatal("expected notification to be marked as read")
	}
}

func TestHTTPNotificationsUnauthorizedV1(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, basePathV1+"/notifications", bytes.NewBufferString(""))
	if err != nil {
		t.Fatalf("notifications unauthorized request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
