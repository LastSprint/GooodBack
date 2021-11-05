export default class OAuth2CodeExchange {
    public code: string
    public provider: string
    public redirect_url: string

    public constructor (code: string, provider: string, redirect_url: string) {
        this.code = code
        this.provider = provider
        this.redirect_url = redirect_url
    }
}
