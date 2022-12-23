package main

// Creates a unique key so we avoid naming collisions in context
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
