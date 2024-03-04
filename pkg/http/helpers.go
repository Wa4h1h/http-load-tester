package http

func getStatusFam(status int) StatusFam {
	switch {
	case status >= 200 && status < 300:

		return success
	case status >= 300 && status < 400:

		return redirection
	case status >= 400 && status < 500:

		return clienterror
	default:

		return servererror
	}
}
