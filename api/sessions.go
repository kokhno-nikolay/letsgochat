package api

func (h *Handler) CheckUserSession(userId int) (bool, string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for key, value := range h.Sessions {
		if value == userId {
			return true, key
		}
	}

	return false, ""
}

func (h *Handler) DeleteSession(token string) error {
	_, ok := h.Sessions[token]
	if ok {
		h.mu.Lock()
		defer h.mu.Unlock()

		if err := h.userRepo.SwitchToInactive(h.Sessions[token]); err != nil {
			return err
		}
		delete(h.Sessions, token)
	}

	return nil
}

func (h *Handler) CheckUserToken(token string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.Sessions[token]
	return ok
}
