package api

func (h *Handler) CheckUserSession(userId int) (bool, string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for key, value := range h.Sessions {
		if value == userId {
			return true, key
		}
	}

	return false, ""
}

func (h *Handler) DeleteSession(token string) {
	_, ok := h.Sessions[token]
	if ok {
		h.mu.Lock()
		defer h.mu.Unlock()
		delete(h.Sessions, token)
	}
}

func (h *Handler) CheckUserToken(token string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.Sessions[token]
	return ok
}
