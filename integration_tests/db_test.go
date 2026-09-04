// Package integration_tests holds integration tests against the real Postgres instance.
//
// They run only when TESTING_DB is set in .env (see Makefile: test-migrate-up
// prepares the schema). Each test runs inside a rolled-back transaction so the
// database is never mutated, and the whole file skips cleanly when no test DB
// is configured — `make test` stays green without Docker.
package integration_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
)

// txStore bundles the SQL store and its enclosing transaction so callers can
// use the repository API while every write is discarded on rollback.
type txStore struct {
	store *repository.SQLStore
	repos repos
	tx    pgx.Tx
}

type repos struct {
	guilds      repository.GuildRepository
	settings    repository.GuildSettingsRepository
	permissions repository.UserPermissionRepository
	audit       repository.AuditLogRepository
}

var uniqueCounter atomic.Int64

func uniqueID(prefix string) string {
	// guilds.id is VARCHAR(20), keep ids short: "g" + ~14 hex chars.
	return fmt.Sprintf("%s%x", prefix, uint64(time.Now().UnixNano())+uint64(10_000+uniqueCounter.Add(1)))
}

// testDBURL derives the test connection string from the app's DB_URL (or, as
// a fallback, GOOSE_DBSTRING) by swapping the database name for TESTING_DB —
// mirroring how the Makefile builds TEST_DB_URL. TESTING_DB may also be a full
// URL on its own. Returns "" when no usable test DB is configured.
func testDBURL() string {
	dbName := os.Getenv("TESTING_DB")
	if dbName == "" {
		return ""
	}
	if strings.Contains(dbName, "://") {
		return dbName
	}

	base := os.Getenv("DB_URL")
	if base == "" {
		base = os.Getenv("GOOSE_DBSTRING")
	}
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	u.Path = "/" + dbName
	return u.String()
}

// newTxStore connects to the TESTING_DB database and starts a transaction.
func newTxStore(t *testing.T) *txStore {
	t.Helper()

	url := testDBURL()
	if url == "" {
		t.Skip("TESTING_DB not set in .env — skipping DB integration tests (test-migrate-up prepares the schema)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Skipf("cannot connect to TESTING_DB (%s) — skipping DB integration tests", err)
	}
	t.Cleanup(func() { _ = conn.Close(context.Background()) })

	tx, err := conn.Begin(context.Background())
	if err != nil {
		t.Skipf("cannot begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })

	store := repository.NewSQLStore(postgres.New(tx))
	return &txStore{
		store: store,
		repos: repos{
			guilds:      repository.NewGuildRepository(store),
			settings:    repository.NewGuildSettingsRepository(store),
			permissions: repository.NewUserPermissionRepository(store),
			audit:       repository.NewAuditLogRepository(store),
		},
		tx: tx,
	}
}

// createGuild inserts a guild and returns its ID.
func createGuild(t *testing.T, ctx context.Context, ts *txStore, name string, tier int32) string {
	t.Helper()
	id := uniqueID("g")
	if _, err := ts.repos.guilds.Create(ctx, dto.CreateGuildParams{ID: id, Name: name, PremiumTier: tier}); err != nil {
		t.Fatalf("create guild: %v", err)
	}
	return id
}

// jsonEqual asserts two byte slices hold semantically-equivalent JSON
// (jsonb columns re-serialize their text form, so raw bytes may differ in
// whitespace and key order).
func jsonEqual(t *testing.T, a, b []byte) bool {
	t.Helper()
	var av, bv any
	if err := json.Unmarshal(a, &av); err != nil {
		t.Errorf("invalid JSON %q: %v", a, err)
		return false
	}
	if err := json.Unmarshal(b, &bv); err != nil {
		t.Errorf("invalid JSON %q: %v", b, err)
		return false
	}
	if !reflect.DeepEqual(av, bv) {
		t.Errorf("JSON not equal: %s vs %s", a, b)
		return false
	}
	return true
}

func TestGuildRepoCRUD(t *testing.T) {
	ctx := context.Background()
	ts := newTxStore(t)
	repo := ts.repos.guilds

	id := createGuild(t, ctx, ts, "Test Guild", 2)

	g, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get guild: %v", err)
	}
	if g.Name != "Test Guild" || g.PremiumTier != 2 || g.CreatedAt.IsZero() {
		t.Errorf("unexpected guild: %+v", g)
	}

	updated, err := repo.Update(ctx, id, dto.UpdateGuildParams{Name: "Renamed", PremiumTier: 5})
	if err != nil {
		t.Fatalf("update guild: %v", err)
	}
	if updated.Name != "Renamed" || updated.PremiumTier != 5 {
		t.Errorf("unexpected updated guild: %+v", updated)
	}

	got, _ := repo.GetByID(ctx, id)
	if got.Name != "Renamed" {
		t.Errorf("get after update = %+v", got)
	}

	all, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list guilds: %v", err)
	}
	found := false
	for _, g := range all {
		if g.ID == id {
			found = true
		}
	}
	if !found {
		t.Errorf("guild %s not present in List", id)
	}

	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("delete guild: %v", err)
	}
	if _, err := repo.GetByID(ctx, id); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows after delete, got %v", err)
	}
}

func TestGuildSettingsRepoUpsertDelete(t *testing.T) {
	ctx := context.Background()
	ts := newTxStore(t)
	repo := ts.repos.settings
	id := createGuild(t, ctx, ts, "Settings Guild", 0)

	upserted, err := repo.Upsert(ctx, dto.UpsertGuildSettingsParams{
		GuildID:               id,
		Prefix:                "-",
		Language:              "en",
		AutoModerationEnabled: true,
		LoggingChannelID:      "123",
		MessageLoggingEnabled: true,
	})
	if err != nil {
		t.Fatalf("upsert settings: %v", err)
	}
	if upserted.Prefix != "-" || upserted.Language != "en" || !upserted.AutoModerationEnabled || upserted.LoggingChannelID != "123" {
		t.Errorf("unexpected upserted settings: %+v", upserted)
	}
	if !upserted.MessageLoggingEnabled {
		t.Errorf("expected MessageLoggingEnabled to persist as true: %+v", upserted)
	}

	got, err := repo.Get(ctx, id)
	if err != nil {
		t.Fatalf("get settings: %v", err)
	}
	if got.Prefix != "-" || got.Language != "en" || !got.MessageLoggingEnabled {
		t.Errorf("unexpected settings: %+v", got)
	}

	updated, err := repo.Update(ctx, id, dto.UpdateGuildSettingsParams{
		Prefix: "!", Language: "de", AutoModerationEnabled: false, LoggingChannelID: "",
		MessageLoggingEnabled: false,
	})
	if err != nil {
		t.Fatalf("update settings: %v", err)
	}
	if updated.Prefix != "!" || updated.Language != "de" || updated.AutoModerationEnabled || updated.LoggingChannelID != "" {
		t.Errorf("unexpected updated settings: %+v", updated)
	}
	if updated.MessageLoggingEnabled {
		t.Errorf("expected MessageLoggingEnabled to update to false: %+v", updated)
	}

	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("delete settings: %v", err)
	}
	if _, err := repo.Get(ctx, id); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows after delete, got %v", err)
	}
}

func TestUserPermissionRepoUpsertListDelete(t *testing.T) {
	ctx := context.Background()
	ts := newTxStore(t)
	repo := ts.repos.permissions
	id := createGuild(t, ctx, ts, "Permission Guild", 0)

	perms := []byte(`["manage_messages","ban"]`)
	up, err := repo.Upsert(ctx, dto.UpsertUserPermissionParams{
		UserID: "user1", GuildID: id, Role: "admin", PermissionsJSON: perms,
	})
	if err != nil {
		t.Fatalf("upsert permission: %v", err)
	}
	if up.Role != "admin" || !jsonEqual(t, up.PermissionsJSON, perms) {
		t.Errorf("unexpected permission: %+v", up)
	}

	got, err := repo.Get(ctx, "user1", id, "admin")
	if err != nil {
		t.Fatalf("get permission: %v", err)
	}
	if got.UserID != "user1" || !jsonEqual(t, got.PermissionsJSON, perms) {
		t.Errorf("unexpected permission: %+v", got)
	}

	// Upsert overwrites rather than duplicates.
	more := []byte(`["kick"]`)
	if _, err := repo.Upsert(ctx, dto.UpsertUserPermissionParams{UserID: "user1", GuildID: id, Role: "admin", PermissionsJSON: more}); err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	if n, err := repo.ListForUser(ctx, "user1", id); err != nil || len(n) != 1 {
		t.Errorf("ListForUser = %d rows, err %v", len(n), err)
	}

	if all, err := repo.ListForGuild(ctx, id); err != nil || len(all) != 1 {
		t.Errorf("ListForGuild = %d rows, err %v", len(all), err)
	}

	if err := repo.Delete(ctx, "user1", id, "admin"); err != nil {
		t.Fatalf("delete permission: %v", err)
	}
	if _, err := repo.Get(ctx, "user1", id, "admin"); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows after delete, got %v", err)
	}
}

func TestAuditLogRepoCRUD(t *testing.T) {
	ctx := context.Background()
	ts := newTxStore(t)
	repo := ts.repos.audit
	id := createGuild(t, ctx, ts, "Audit Guild", 0)

	created, err := repo.Create(ctx, dto.CreateAuditLogParams{
		GuildID:  id,
		Action:   "MESSAGE_DELETE",
		ActorID:  "user9",
		TargetID: "channel1",
		Metadata: []byte(`{"message_id":"123"}`),
	})
	if err != nil {
		t.Fatalf("create audit log: %v", err)
	}
	if created.ID == 0 || created.Action != "MESSAGE_DELETE" || created.CreatedAt.IsZero() {
		t.Errorf("unexpected audit log: %+v", created)
	}
	if created.TargetID != "channel1" {
		t.Errorf("target id = %q", created.TargetID)
	}

	count, err := repo.CountForGuild(ctx, id)
	if err != nil || count != 1 {
		t.Fatalf("count = %d, err %v", count, err)
	}

	_, _ = repo.Create(ctx, dto.CreateAuditLogParams{
		GuildID:  id,
		Action:   "MESSAGE_CREATE",
		ActorID:  "user9",
		TargetID: "channel1",
		Metadata: []byte(`{"message_id":"456"}`),
	})

	if count, _ := repo.CountForGuild(ctx, id); count != 2 {
		t.Errorf("count after second insert = %d", count)
	}

	forGuild, err := repo.ListForGuild(ctx, id, 10, 0)
	if err != nil || len(forGuild) != 2 {
		t.Fatalf("ListForGuild = %d rows, err %v", len(forGuild), err)
	}
	// Both rows were inserted in the same transaction, so created_at ties —
	// only the set of actions is deterministic.
	actions := map[string]bool{}
	for _, el := range forGuild {
		actions[el.Action] = true
	}
	if !actions["MESSAGE_CREATE"] || !actions["MESSAGE_DELETE"] {
		t.Errorf("expected both action types, got %v", actions)
	}
	foundMeta := false
	for _, el := range forGuild {
		if el.Action == "MESSAGE_CREATE" && jsonEqual(t, el.Metadata, []byte(`{"message_id":"456"}`)) {
			foundMeta = true
		}
	}
	if !foundMeta {
		t.Errorf("missing MESSAGE_CREATE entry with message_id 456: %+v", forGuild)
	}

	byActor, err := repo.ListByActor(ctx, id, "user9", 10, 0)
	if err != nil || len(byActor) != 2 {
		t.Errorf("ListByActor = %d rows, err %v", len(byActor), err)
	}
	if byActor, _ := repo.ListByActor(ctx, id, "nobody", 10, 0); len(byActor) != 0 {
		t.Errorf("ListByActor(nobody) = %d rows", len(byActor))
	}

	// DeleteBefore: a historic cutoff keeps new rows; a future cutoff clears them.
	if err := repo.DeleteBefore(ctx, time.Now().Add(-time.Hour)); err != nil {
		t.Fatalf("delete before (past): %v", err)
	}
	if count, _ := repo.CountForGuild(ctx, id); count != 2 {
		t.Errorf("past cutoff should not delete, count = %d", count)
	}
	if err := repo.DeleteBefore(ctx, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("delete before (future): %v", err)
	}
	if count, _ := repo.CountForGuild(ctx, id); count != 0 {
		t.Errorf("future cutoff should delete all, count = %d", count)
	}
}

// TestAuditLogRepoDeleteMessageLogsBefore verifies the retention prune removes
// only MESSAGE_* rows and leaves moderation entries in place.
func TestAuditLogRepoDeleteMessageLogsBefore(t *testing.T) {
	ctx := context.Background()
	ts := newTxStore(t)
	repo := ts.repos.audit
	id := createGuild(t, ctx, ts, "Retention Guild", 0)

	mustCreate := func(action string) {
		if _, err := repo.Create(ctx, dto.CreateAuditLogParams{
			GuildID: id, Action: action, ActorID: "user1", TargetID: "chan1",
			Metadata: []byte(`{}`),
		}); err != nil {
			t.Fatalf("create %s: %v", action, err)
		}
	}
	mustCreate("MESSAGE_CREATE")
	mustCreate("MESSAGE_DELETE")
	mustCreate("BAN") // moderation action — must survive the prune

	// A future cutoff would match every row by time; only MESSAGE_* should go.
	if err := repo.DeleteMessageLogsBefore(ctx, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("delete message logs: %v", err)
	}

	remaining, err := repo.ListForGuild(ctx, id, 10, 0)
	if err != nil {
		t.Fatalf("list after prune: %v", err)
	}
	if len(remaining) != 1 {
		t.Fatalf("expected 1 row after pruning message logs, got %d: %+v", len(remaining), remaining)
	}
	if remaining[0].Action != "BAN" {
		t.Errorf("expected the BAN moderation entry to survive, got %q", remaining[0].Action)
	}
}
