export const SERVER_BASE_URL: string = process.env.VUE_APP_SERVER_BASE_URL
export const WEB_BASE_URL: string = process.env.VUE_APP_WEB_BASE_URL
export const GOOGLE_CLIENT_ID: string = process.env.VUE_APP_GOOGLE_CLIENT_ID

export const REDIRECT_PATH = "/auth-redirect"
export const REDIRECT_URL: string = WEB_BASE_URL + REDIRECT_PATH
export const USE_TLS: boolean = process.env.VUE_APP_USE_TLS
export const TURN_LOGGER: boolean = process.env.VUE_APP_TURN_LOGGER
