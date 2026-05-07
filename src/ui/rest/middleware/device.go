package middleware

import (
	"net/url"
	"strings"

	"whatsapp-bot/config"
	"whatsapp-bot/infrastructure/whatsapp"
	"whatsapp-bot/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

const DeviceIDHeader = "X-Device-Id"

// deviceScopedPrefixes lists the path prefixes that require a device context.
// Any request whose path (after stripping AppBasePath) does NOT start with one of
// these prefixes is allowed to pass through so that Fiber can return its natural
// 404 — rather than a misleading 400 DEVICE_ID_REQUIRED response.
var deviceScopedPrefixes = []string{
	"/app/",
	"/send/",
	"/user/",
	"/chat/",
	"/message/",
	"/group/",
	"/newsletter/",
	"/ws",
}

// isDeviceScopedPath returns true if path requires a device context.
func isDeviceScopedPath(path string) bool {
	// Strip AppBasePath prefix so comparisons work regardless of deployment prefix.
	trimmed := path
	if bp := config.AppBasePath; bp != "" {
		trimmed = strings.TrimPrefix(path, bp)
	}
	if trimmed == "" {
		trimmed = "/"
	}
	for _, prefix := range deviceScopedPrefixes {
		if trimmed == strings.TrimSuffix(prefix, "/") || strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

// DeviceMiddleware fetches a device instance by header (preferred), path param, or query param
// and injects it into the context. It falls back to the default/only device for single-device mode.
func DeviceMiddleware(dm *whatsapp.DeviceManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := strings.TrimSpace(c.Path())

		// Pass through root and any path that is not device-scoped.
		// Unknown paths will reach Fiber's natural 404 handler instead of
		// getting a misleading 400 DEVICE_ID_REQUIRED response.
		if path == "/" || path == "" || path == config.AppBasePath || path == config.AppBasePath+"/" {
			return c.Next()
		}
		if !isDeviceScopedPath(path) {
			return c.Next()
		}

		if dm == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(utils.ResponseData{
				Status:  fiber.StatusServiceUnavailable,
				Code:    "DEVICE_MANAGER_UNAVAILABLE",
				Message: "Device manager is not initialized",
				Results: nil,
			})
		}

		deviceID := strings.TrimSpace(c.Get(DeviceIDHeader))
		// URL-decode the header value to support non-ASCII characters
		if decoded, err := url.QueryUnescape(deviceID); err == nil {
			deviceID = decoded
		}
		if deviceID == "" {
			deviceID = strings.TrimSpace(c.Query("device_id"))
		}

		instance, resolvedID, err := dm.ResolveDevice(deviceID)
		if err != nil {
			// ResolveDevice returns an ID when provided but missing; use it for payload clarity.
			if resolvedID != "" || strings.TrimSpace(deviceID) != "" {
				return c.Status(fiber.StatusNotFound).JSON(utils.ResponseData{
					Status:  fiber.StatusNotFound,
					Code:    "DEVICE_NOT_FOUND",
					Message: "device not found; create a device first from /api/devices or provide a valid X-Device-Id",
					Results: map[string]string{"device_id": resolvedID},
				})
			}

			return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseData{
				Status:  fiber.StatusBadRequest,
				Code:    "DEVICE_ID_REQUIRED",
				Message: "device_id is required via X-Device-Id header or device_id query",
				Results: nil,
			})
		}

		c.Locals("device_id", resolvedID)
		c.Locals("device", instance)
		c.SetUserContext(whatsapp.ContextWithDevice(c.UserContext(), instance))
		return c.Next()
	}
}
