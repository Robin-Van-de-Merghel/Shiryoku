# Shiryoku-routers

This layer is the only one (for now) exposed to the internet. It receives the requests, check for authentication, and forward it to the logic layer.

# Global Architecture

The current architecture for the API is:

1. `/api/auth/*` for the auth process
2. `/api/config/*` to fetch the current config
3. `/api/agents/*` for the agents to push/pull data, or administrators to get information about them
4. `/api/modules/*` for module specific interactions
5. `/api/widgets/*` for widgets specific interactions

## Users interactions

> [!NOTE]
> This part is in construction. First we need a working product before adding auth.

## Configurations

Some of the configuration parts may be fetched by agents or users, such as:

1. Version(s)
2. Scopes

> [!NOTE]
> To be completed.

## Agents

Agents are the one exploring the internet and pushing data. Special routes are dedicated to them:

1. `POST /api/agents/ping`: to let the server know they are alive
2. `GET /api/agents/tasks`: to fetch tasks to run 
3. `GET /api/agents/configure`: to fetch configuration for some tools
4. `POST /api/agents/modules/{XXX}/upload`: to upload data for a specific module
5. `PATCH /api/agents/modules/{XXX}/upload`: to modify data already sent

> [!IMPORTANT]
> As this part is in construction, it might change a lot!

## Modules

This sub-section of the API is dedicated to dashboard to fetch data about specific modules.

> [!NOTE]
> See if we might keep only one generic endpoint, or multiple categorized by module.

# Dependency injection

To inject data, I used [Alex Edwards](https://www.alexedwards.net/blog/organising-database-access)'s guidelines, as such:

```go
// Handler takes the injected values as parameters
// Return what a gin.Handler function requires: a function taking a context as an input
func YourHandler(ASqlDB *sql.DB) func(c *gin.Context) {
	// Wraps to inject dependency (ASqlDB)
    return func(c *gin.Context) {
        // Your classical handler
        ...
    }
}
```

> [!NOTE]
> `HandlerFunc` is just a redefinition of `func(*gin.Context)` (see [here](https://pkg.go.dev/github.com/gin-gonic/gin#HandlerFunc))

Then, you can call it in [`routers.go`](./routers.go) via:

```go
router.GET(
    "/test", 
    // Call the handler to inject a value
    // (normally, you would just give a reference to the handler)
    YourHandler(ASqlDB),
)
```
