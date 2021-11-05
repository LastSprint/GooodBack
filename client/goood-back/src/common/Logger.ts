import { TURN_LOGGER } from "./GlobalConstants"

export default class Logger {
    static log (...data: any): void {

        if (TURN_LOGGER) {
            console.log(data)
        }
    }
}
