package models

import (
"reflect"
"strings"
"testing"
"time"
)

func TestRoleTableName(t *testing.T) {
if table := (Role{}).TableName(); table != "roles" {
t.Fatalf("expected roles table, got %s", table)
}
}

func TestUserRoleTableName(t *testing.T) {
if table := (UserRole{}).TableName(); table != "user_roles" {
t.Fatalf("expected user_roles table, got %s", table)
}
}

func TestRefreshTokenTableName(t *testing.T) {
if table := (RefreshToken{}).TableName(); table != "refresh_tokens" {
t.Fatalf("expected refresh_tokens table, got %s", table)
}
}

func TestRoleShape(t *testing.T) {
typeOfRole := reflect.TypeOf(Role{})

idField, ok := typeOfRole.FieldByName("ID")
if !ok || idField.Type != reflect.TypeOf(uint(0)) {
t.Fatalf("Role.ID must be uint")
}
if !strings.Contains(idField.Tag.Get("gorm"), "primaryKey") {
t.Fatalf("Role.ID must define primaryKey gorm tag")
}

nameField, ok := typeOfRole.FieldByName("Name")
if !ok || nameField.Type != reflect.TypeOf("") {
t.Fatalf("Role.Name must be string")
}
gormTag := nameField.Tag.Get("gorm")
if !strings.Contains(gormTag, "not null") || !strings.Contains(gormTag, "unique") {
t.Fatalf("Role.Name must be not null and unique")
}
}

func TestUserRoleShape(t *testing.T) {
typeOfUserRole := reflect.TypeOf(UserRole{})

idField, ok := typeOfUserRole.FieldByName("ID")
if !ok || idField.Type != reflect.TypeOf(uint(0)) {
t.Fatalf("UserRole.ID must be uint")
}
if !strings.Contains(idField.Tag.Get("gorm"), "primaryKey") {
t.Fatalf("UserRole.ID must define primaryKey gorm tag")
}

userIDField, ok := typeOfUserRole.FieldByName("UserID")
if !ok || userIDField.Type != reflect.TypeOf(uint(0)) {
t.Fatalf("UserRole.UserID must be uint")
}
if !strings.Contains(userIDField.Tag.Get("gorm"), "not null") {
t.Fatalf("UserRole.UserID must be not null")
}

roleIDField, ok := typeOfUserRole.FieldByName("RoleID")
if !ok || roleIDField.Type != reflect.TypeOf(uint(0)) {
t.Fatalf("UserRole.RoleID must be uint")
}
if !strings.Contains(roleIDField.Tag.Get("gorm"), "not null") {
t.Fatalf("UserRole.RoleID must be not null")
}
}

func TestRefreshTokenShape(t *testing.T) {
typeOfRefreshToken := reflect.TypeOf(RefreshToken{})

idField, ok := typeOfRefreshToken.FieldByName("ID")
if !ok || idField.Type != reflect.TypeOf(uint(0)) {
t.Fatalf("RefreshToken.ID must be uint")
}
if !strings.Contains(idField.Tag.Get("gorm"), "primaryKey") {
t.Fatalf("RefreshToken.ID must define primaryKey gorm tag")
}

tokenHashField, ok := typeOfRefreshToken.FieldByName("TokenHash")
if !ok || tokenHashField.Type != reflect.TypeOf("") {
t.Fatalf("RefreshToken.TokenHash must be string")
}
tokenHashTag := tokenHashField.Tag.Get("gorm")
if !strings.Contains(tokenHashTag, "not null") || !strings.Contains(tokenHashTag, "unique") {
t.Fatalf("RefreshToken.TokenHash must be not null and unique")
}

userIDField, ok := typeOfRefreshToken.FieldByName("UserID")
if !ok || userIDField.Type != reflect.TypeOf(uint(0)) {
t.Fatalf("RefreshToken.UserID must be uint")
}
if !strings.Contains(userIDField.Tag.Get("gorm"), "not null") {
t.Fatalf("RefreshToken.UserID must be not null")
}

expiresAtField, ok := typeOfRefreshToken.FieldByName("ExpiresAt")
if !ok || expiresAtField.Type != reflect.TypeOf(time.Time{}) {
t.Fatalf("RefreshToken.ExpiresAt must be time.Time")
}
if !strings.Contains(expiresAtField.Tag.Get("gorm"), "not null") {
t.Fatalf("RefreshToken.ExpiresAt must be not null")
}

revokedAtField, ok := typeOfRefreshToken.FieldByName("RevokedAt")
if !ok || revokedAtField.Type != reflect.TypeOf((*time.Time)(nil)) {
t.Fatalf("RefreshToken.RevokedAt must be *time.Time")
}

createdAtField, ok := typeOfRefreshToken.FieldByName("CreatedAt")
if !ok || createdAtField.Type != reflect.TypeOf(time.Time{}) {
t.Fatalf("RefreshToken.CreatedAt must be time.Time")
}
}
