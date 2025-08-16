# OAuth Credentials File Support for Anthropic

*Add support for reading OAuth credentials from a pre-existing credentials.json file as an alternative to API key authentication, while preserving the ability to use different authentication methods across sessions.*

## Requirements

- Add `--auth-oauth` command-line flag to trigger OAuth credentials file authentication for current session
- When flag is used, read OAuth credentials from `credentials.json` file in working directory
- Parse JSON structure to extract `accessToken`, `refreshToken`, and `expiresAt` from `claudeAiOauth` section
- Check token expiration and automatically refresh if expired using refresh token
- Preserve existing session isolation - different terminals can use different auth methods
- Fall back gracefully to API key authentication if credentials file missing/invalid or OAuth fails
- Use OAuth access token for API requests when available and valid for current session
- Handle token refresh transparently using Anthropic's token refresh endpoint (requires CLIENT_ID)

## Rules

- OAuth credentials file authentication is opt-in via command-line flag only per session
- API key authentication remains the default method
- Environment variables (AGENT_API_KEY, ANTHROPIC_API_KEY) always take precedence when set
- Session-specific auth choice should not affect other running sessions
- Credentials file (`credentials.json`) must exist in working directory when `--auth-oauth` flag is used
- Failed OAuth attempts should not prevent fallback to API key auth
- Token refresh should be transparent to the user and update the credentials file (requires ANTHROPIC_CLIENT_ID)
- Only implement for Anthropic provider initially

## Domain

```go
// OAuth credentials from credentials.json file
type ClaudeAiOauth struct {
    AccessToken      string   `json:"accessToken"`
    RefreshToken     string   `json:"refreshToken"`
    ExpiresAt        int64    `json:"expiresAt"`        // Unix timestamp in milliseconds
    Scopes           []string `json:"scopes"`
    SubscriptionType string   `json:"subscriptionType"`
}

// Root credentials file structure
type CredentialsFile struct {
    ClaudeAiOauth ClaudeAiOauth `json:"claudeAiOauth"`
}

// Authentication method priority
type AuthConfig struct {
    UseOAuth          bool   // Set by --auth-oauth flag
    APIKey            string // From environment variables
    CredentialsFile   *CredentialsFile // Loaded from credentials.json if available
}
```

## Extra Considerations

- Authentication precedence: Environment variables > OAuth credentials file > failure
- `credentials.json` file must be in working directory, not `.agent/` subdirectory
- Session logging should indicate which auth method was used
- Handle concurrent access to credentials file if multiple sessions refresh simultaneously
- File locking during token refresh to prevent corruption
- Preserve all other fields in credentials.json when updating tokens
- Non-OAuth sessions should ignore OAuth credential file issues
- **REQUIRED**: Token refresh requires ANTHROPIC_CLIENT_ID environment variable - must be obtained from Anthropic's OAuth configuration
- If ANTHROPIC_CLIENT_ID is not available, skip token refresh and let tokens expire (graceful degradation)

## Testing Considerations

- Test session isolation: OAuth in one terminal, API key in another
- Test environment variable precedence over OAuth tokens
- Unit tests for authentication method selection logic
- Mock HTTP client tests for token exchange endpoints
- Test token storage and retrieval from file system
- Test concurrent token refresh from multiple sessions
- Test graceful fallback scenarios

## Implementation Notes

- Modify `main()` to check `--auth-oauth` flag and set session-specific auth preference
- Authentication selection logic:
  1. If `AGENT_API_KEY` or `ANTHROPIC_API_KEY` set → use API key
  2. Else if `--auth-oauth` flag set → attempt OAuth credentials file (with fallback to API key if fails)
  3. Else → use API key authentication (existing behavior)
- Read credentials from `credentials.json` in working directory
- Token expiration check: compare `expiresAt` (milliseconds) with current time
- Add file locking for concurrent token refresh operations
- Use Anthropic's token refresh endpoint: `https://console.anthropic.com/v1/oauth/token` (requires ANTHROPIC_CLIENT_ID environment variable)
- Update credentials file in-place when tokens are refreshed, preserving structure

## Specification by Example

**Command usage:**
```bash
./agent --auth-oauth
```

**OAuth credentials file authentication:**
```
Loading OAuth credentials from credentials.json...
Using OAuth authentication with Anthropic
```

**Session 1 with OAuth:**
```bash
ANTHROPIC_API_KEY="" ./agent --auth-oauth
# Reads credentials.json, uses OAuth access token
```

**Session 2 with API key (simultaneous):**
```bash
ANTHROPIC_API_KEY="sk-123" ./agent
# Uses API key, ignores credentials.json
```

**Session 3 reusing OAuth:**
```bash
ANTHROPIC_API_KEY="" ./agent --auth-oauth
# Reuses credentials from credentials.json
```

**Authentication precedence example:**
```bash
ANTHROPIC_API_KEY="sk-123" ./agent --auth-oauth
# Still uses API key despite --auth-oauth flag
```

**Credentials file structure (`credentials.json`):**
```json
{
  "claudeAiOauth": {
    "accessToken": "sk-ant-oREDACTED8xtkTYZLaAljUYvfHPewz5Eh-Q-5ZBCngAA",
    "refreshToken": "sk-ant-ort01-vLIsai0Wc_1i-rGQEREDACTED1rLhhk9sy-hcU5L3Jw24GmjpWcT9CGNA-ZJOmrAAA",
    "expiresAt": 1754921476460,
    "scopes": ["user:inference", "user:profile"],
    "subscriptionType": "pro"
  }
}
```

## Reference Implementation Details

### Token Refresh (from TypeScript reference)

The token refresh logic requires CLIENT_ID as shown in the original OAuth implementation:

```typescript
// Token refresh endpoint and logic - requires CLIENT_ID
const response = await fetch(
  "https://console.anthropic.com/v1/oauth/token",
  {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      grant_type: "refresh_token",
      refresh_token: info.refresh,
      client_id: CLIENT_ID,  // ← REQUIRED: Must be provided via ANTHROPIC_CLIENT_ID env var
    }),
  },
)
```

### Credentials File Example

Based on the provided sanitized example:

```json
{
  "claudeAiOauth": {
    "accessToken": "sk-ant-oREDACTED8xtkTYZLaAljUYvfHPewz5Eh-Q-5ZBCngAA",
    "refreshToken": "sk-ant-ort01-vLIsai0Wc_1i-rGQEREDACTED1rLhhk9sy-hcU5L3Jw24GmjpWcT9CGNA-ZJOmrAAA",
    "expiresAt": 1754921476460,
    "scopes": ["user:inference", "user:profile"],
    "subscriptionType": "pro"
  }
}
```

### Key Implementation Details:

1. **File Reading**: Read `credentials.json` from working directory
2. **Expiration Check**: Compare `expiresAt` (Unix timestamp in milliseconds) with current time
3. **Token Refresh**: Use refresh token with Anthropic's endpoint when expired
4. **File Update**: Update credentials.json with new tokens after refresh, preserving structure
5. **Concurrent Access**: Use file locking to prevent corruption during updates

## Verification

- [ ] `--auth-oauth` flag reads credentials from `credentials.json` when no API key environment variable set
- [ ] Environment variables always take precedence over OAuth credentials file
- [ ] Multiple terminal sessions can use different authentication methods simultaneously
- [ ] OAuth credentials are reused across sessions when `--auth-oauth` is specified
- [ ] Failed OAuth in one session doesn't affect API key usage in other sessions
- [ ] Token refresh works correctly with concurrent sessions and updates credentials.json
- [ ] Session logs indicate which authentication method was used
- [ ] Fallback to API key works when credentials file is missing/invalid or refresh fails
- [ ] Expiration check correctly handles Unix timestamp in milliseconds format
- [ ] Token refresh request format matches Anthropic's API requirements
- [ ] ANTHROPIC_CLIENT_ID environment variable is used in refresh requests when available
- [ ] Token refresh is skipped gracefully when ANTHROPIC_CLIENT_ID is not set
- [ ] File locking prevents corruption during concurrent credential updates
- [ ] All fields in credentials.json are preserved during token refresh

## Next Steps

After creating the story:
1. Save to `.agent/stories/oauth-authentication/user-story.md`