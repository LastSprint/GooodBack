import * as global from '../common/GlobalConstants'
import axios from 'axios'
import { Feedback } from './entries/Feedback'
import { NewFeedback } from './entries/NewFeedback'

export default class FeedbackService {

    loadServices (): Promise<[Feedback]> {
        return axios.get(global.SERVER_BASE_URL + "/feedback", { withCredentials: true }).then((res) => {
            return res.data
        })
    }

    send (fedback: NewFeedback): Promise<void> {
        return axios.post(global.SERVER_BASE_URL + "/feedback", fedback)
    }
}
