export class Feedback {
    id?: string
    message!: string
    creation_date!: string
    type!: number

    public emoji (): string {
        switch (this.type) {
            case 0: return '👍'
            case 1: return '👎'
            case 2: return '🔥'
            case 3: return '🤬'
        }

        return '😨'
    }
}
