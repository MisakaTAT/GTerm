package sftp

import (
	"bufio"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
	"golang.org/x/crypto/ssh"
)

type FileOwnerGroup struct {
	owner string
	group string
}

type PermissionsCache struct {
	mu            sync.Mutex
	users         map[uint32]string
	groups        map[uint32]string
	fileInfoCache map[string]FileOwnerGroup
}

func NewPermissionsCache() *PermissionsCache {
	return &PermissionsCache{
		users:         make(map[uint32]string),
		groups:        make(map[uint32]string),
		fileInfoCache: make(map[string]FileOwnerGroup),
	}
}

func (c *PermissionsCache) GetUsername(uid uint32) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name, ok := c.users[uid]; ok {
		return name
	}
	return strconv.Itoa(int(uid))
}

func (c *PermissionsCache) GetGroupName(gid uint32) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name, ok := c.groups[gid]; ok {
		return name
	}
	return strconv.Itoa(int(gid))
}

func (c *PermissionsCache) preloadPermissions(conn *ssh.Client, execAdapter *exec.Adapter) {
	if conn == nil {
		slog.Warn("Cannot preload permissions: SSH client is nil")
		return
	}

	slog.Info("Preloading users and groups information")
	users, groups := c.fetchAllUsersAndGroups(execAdapter)

	c.mu.Lock()
	defer c.mu.Unlock()

	for uid, name := range users {
		c.users[uid] = name
	}
	for gid, name := range groups {
		c.groups[gid] = name
	}

	slog.Info("Preloaded %d users and %d groups", len(users), len(groups))
}

func (c *PermissionsCache) fetchAllUsersAndGroups(execAdapter *exec.Adapter) (map[uint32]string, map[uint32]string) {
	users := make(map[uint32]string)
	groups := make(map[uint32]string)

	slog.Debug("Executing command to fetch users and groups")
	result := execAdapter.Run("(echo '===USERS==='; getent passwd; echo '===GROUPS==='; getent group)")
	if !result.Success() {
		slog.Warn("Failed to fetch users and groups using getent: %v", result.Error())
		slog.Info("Trying fallback method to read users and groups")
		return c.fetchUsersAndGroupsFallback(execAdapter)
	}

	scanner := bufio.NewScanner(strings.NewReader(result.StdOut()))
	mode := ""

	for scanner.Scan() {
		line := scanner.Text()
		if line == "===USERS===" {
			mode = "users"
			continue
		} else if line == "===GROUPS===" {
			mode = "groups"
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}
		if mode == "users" {
			uid, err := strconv.ParseUint(parts[2], 10, 32)
			if err == nil {
				users[uint32(uid)] = parts[0]
			}
		} else if mode == "groups" {
			gid, err := strconv.ParseUint(parts[2], 10, 32)
			if err == nil {
				groups[uint32(gid)] = parts[0]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error scanning output: %v", err)
	}

	return users, groups
}

func (c *PermissionsCache) fetchUsersAndGroupsFallback(execAdapter *exec.Adapter) (map[uint32]string, map[uint32]string) {
	users := make(map[uint32]string)
	groups := make(map[uint32]string)

	passwdResult := execAdapter.Run("cat /etc/passwd")
	if passwdResult.Success() {
		slog.Debug("Reading users from /etc/passwd")
		scanner := bufio.NewScanner(strings.NewReader(passwdResult.StdOut()))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) < 3 {
				continue
			}
			uid, err := strconv.ParseUint(parts[2], 10, 32)
			if err == nil {
				users[uint32(uid)] = parts[0]
			}
		}
		if err := scanner.Err(); err != nil {
			slog.Error("Error scanning passwd output: %v", err)
		}
	} else {
		slog.Warn("Failed to read /etc/passwd: %v", passwdResult.Error())
	}

	groupResult := execAdapter.Run("cat /etc/group")
	if groupResult.Success() {
		slog.Debug("Reading groups from /etc/group")
		scanner := bufio.NewScanner(strings.NewReader(groupResult.StdOut()))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) < 3 {
				continue
			}
			gid, err := strconv.ParseUint(parts[2], 10, 32)
			if err == nil {
				groups[uint32(gid)] = parts[0]
			}
		}
		if err := scanner.Err(); err != nil {
			slog.Error("Error scanning group output: %v", err)
		}
	} else {
		slog.Warn("Failed to read /etc/group: %v", groupResult.Error())
	}

	return users, groups
}
