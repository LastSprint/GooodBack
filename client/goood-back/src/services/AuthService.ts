import axios from 'axios'
import OAuth2CodeExchange from './entries/OAuth2CodeExchange'
import Tokens from './entries/Tokens'
import * as global from '../common/GlobalConstants'
import Vue from 'vue'
import Logger from '@/common/Logger'

export default class AuthService {
    private _authCookieName = 'Authorization'
    private _refreshTokenCookieName = 'Refreshing'

    isAuthorizaed (): boolean {
        return this._getCookie(this._authCookieName) != null
    }

    authorize (code: string): Promise<void> {
        const model = new OAuth2CodeExchange(code, 'google', global.REDIRECT_URL)
        return axios.post<Tokens>(`${global.SERVER_BASE_URL}/auth/thirparty/code`, model, { withCredentials: true }).then((res) => {
            Logger.log(document.cookie)
            Logger.log(res.headers)
        })
    }
    
    private _getCookie (name: string): string|null {
        const nameLenPlus = (name.length + 1)

        return document.cookie
            .split(';')
            .map(c => c.trim())
            .filter(cookie => {
                return cookie.substring(0, nameLenPlus) === `${name}=`
            })
            .map(cookie => {
                return decodeURIComponent(cookie.substring(nameLenPlus))
            })[0]
    }
}
