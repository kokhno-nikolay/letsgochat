package api

func (h *Handler) CheckUserSession(userId int) (bool, string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for key, value := range h.sessions {
		if value == userId {
			return true, key
		}
	}

	return false, ""
}

func (h *Handler) DeleteSession(token string) error {
	_, ok := h.sessions[token]
	if ok {
		h.mu.Lock()
		defer h.mu.Unlock()

		if err := h.userRepo.SwitchToInactive(h.sessions[token]); err != nil {
			return err
		}
		delete(h.sessions, token)
	}

	return nil
}

func (h *Handler) CheckUserToken(token string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.sessions[token]
	return ok
}
