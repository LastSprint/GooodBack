export class NewFeedback {
    message: string
    type: number
    target: string

    constructor (message: string, type: number, target: string) {
        this.message = message
        this.type = type
        this.target = target
    }
}
