<template>
    <div class>
          <vs-alert relief v-model="isError" color="#d90429" closable>
            <template #title>
            Error occured
            </template>
            <p>{{ errorMessage }}</p>
        </vs-alert>
        <vs-alert relief v-model="isSuccess" color="#74c69d" closable>
            <p>Feedback was sent 📫</p>
        </vs-alert>

        <vs-row style="margin-top:10%">
            <vs-col offset="5">
                <vs-input
                    placeholder="email"
                    v-model="email"
                    style="color:#e9ecef; font-size: 14pt;">
                    <template #icon>
                        <img src="../assets/user_ic.png" width="24"/>
                    </template>
                </vs-input>
            </vs-col>
        </vs-row>

        <vs-row style="margin-top:24px" justify="center">
            <vs-col offset="5">
                <vs-select style="color:#e9ecef;margin-top:24" v-model="value">
                    <vs-option label="👍" value=0>
                    👍
                    </vs-option>
                    <vs-option label="👎" value=1>
                    👎
                    </vs-option>
                    <vs-option label="🔥" value=2>
                    🔥
                    </vs-option>
                    <vs-option label="🤬" value=3>
                    🤬
                    </vs-option>
                </vs-select>
            </vs-col>
        </vs-row> 
        <vs-row style="margin-top:24px" justify="center">
            <vs-col offset="5">
                <vs-card style="max-width:100000px">
                    <template #text>
                        <textarea style="background:#212529;border:0px" cols=32 rows=12 class="textArea" v-model="message"></textarea>
                    </template>
                </vs-card>
            </vs-col>
        </vs-row>
            <vs-row style="margin-top:24px" justify="center">
                <vs-col offset="5">
                    <vs-button
                        v-bind:loading="formWasSent"
                        success
                        size="xl"
                        style="width:160px"
                        @click="send"
                        flat>
                            Send
                    </vs-button>
                </vs-col>
        </vs-row>
    </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { Component, Prop } from 'vue-property-decorator'
import { NewFeedback } from '../services/entries/NewFeedback'
import FeedbackService from '../services/FeedbackService'

@Component
export default class SendFeedback extends Vue {

    isLoading = true
    types: number[] = [0, 1, 2, 3]
    value = 0
    email = ""
    message = ""

    formWasSent = false
    isError = false
    errorMessage = ""

    isSuccess = true

    created () {
        this.$vs.setColor('gray-2', '#212529')
        this.$vs.setColor('gray-1', '#212529')
        this.$vs.setColor('gray-3', '#212529')
        this.$vs.setColor('background', '#212529')
        this.$vs.setColor('text', '#e9ecef')
    }

    onAlertClick (): void {
        this.isError = false
    }

    send (): void {
        
        const regexp = new RegExp("[a-z0-9-_\\.]+@surfstudio.ru")
        const result = regexp.exec(this.email)

        if (result?.length === 0 || result == null) {
            this.email = ""
            this.errorMessage = "Email field was incorrect. Please, make sure that the email has '@surfstudio.ru' domain"
            this.isSuccess = false
            this.isError = true
            return
        }

        if (this.message.trim().length === 0) {
            this.errorMessage = "Sorry but you have to leave a message :)"
            this.isSuccess = false
            this.isError = true
            return
        }

        this.isSuccess = false
        this.isError = false

        const srv = new FeedbackService()
        this.formWasSent = true
        srv.send(new NewFeedback(this.message, parseInt(`${this.value}`), this.email)).then((response) => {
            window.alert('Success')
        }).catch((error) => {
                if (error.response) {
                    this.errorMessage = `Error: ${error.response.data}`
                    this.isSuccess = true
                } else if (error.request) {
                    console.log(error.request)
                    this.errorMessage = `Can't connect to server ☠️`
                } else {
                    this.errorMessage = `Something went wrong 😨`
                }

                this.isSuccess = false
                this.isError = true
        }).finally(() => {
            this.formWasSent = false
        })
    }
}
</script>

<style scoped>

p{
 font-family: 'Comic Sans MS', 'Comic Sans', cursive;
 font-size: 14pt;
 color: #e9ecef;
}

.vs-card {
    max-width: 100%;
    background: #212529; 
}

.reaction {
    border:#6c757d;
    border-left-style: solid;
    border-width: 1px;
}

.time {
    border:#e9ecef;
    border-top-style: solid;
    border-width: 1px;
}

.textArea {
    margin: 6px;
    resize: none;
     font-family: 'Comic Sans MS', 'Comic Sans', cursive;
    font-size: 14pt;
    color: #e9ecef;
}

.textArea:focus {
    border:0px;
    outline: none !important;
}

.pageCenter {
    position: fixed;
    top: 50%;
    -webkit-transform: translate(0, -50%);
    transform: translate(0, -50%);
}

</style>
