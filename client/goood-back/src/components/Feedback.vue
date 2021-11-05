<template>
    <div>
        <vs-navbar center-collapsed>
            <template #left>
                <img src="../assets/login_logo.png" width="32px" height="32px" style="margin-top:6px;margin-left:6px;margin-bottom:6px;margin-right:6px">
                <p> GoodBack </p>
                <p style="margin-left: 48px"> {{ typeToEmoji(0) }} {{ countOf(0) }} </p>
                <p style="margin-left: 24px"> {{ typeToEmoji(1) }} {{ countOf(1) }} </p>
                <p style="margin-left: 24px"> {{ typeToEmoji(2) }} {{ countOf(2) }} </p>
                <p style="margin-left: 24px"> {{ typeToEmoji(3) }} {{ countOf(3) }} </p>
            </template>
        </vs-navbar>
        <div v-if="isLoading" style="margin-top:65px">
            LOADING ...
        </div>
        <div v-else style="margin-top:90px">
            <div v-if="feedbacks.length != 0">
                <vs-row justify="center" v-for="item in feedbacks" v-bind:key="item.id" style="margin-top:12px">
                    <vs-col w=8 >
                        <div class="vs-card">
                            <vs-row>
                                <vs-col w="1">
                                    <h2>{{ typeToEmoji(item.type) }}</h2>
                                </vs-col>
                                <vs-col w="11" class="reaction">
                                    <vs-row justify="center">
                                        <vs-col w="12">
                                            <p> {{ item.message }} </p>
                                            <p class="date-p"> {{ formatData(item.creation_date) }} </p>
                                        </vs-col>
                                    </vs-row>
                                </vs-col>
                            </vs-row>
                        </div>
                    </vs-col>
                </vs-row>
            </div>
            <div v-else>
                EMPTY
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { Component, Prop } from 'vue-property-decorator'
import { Feedback } from '../services/entries/Feedback'
import FeedbackService from '../services/FeedbackService'

@Component
export default class FeedbackCmp extends Vue {

    isLoading = true
    feedbacks: Feedback[] = []

    formatData (val: string): string {
        const dateTimeFormat = new Intl.DateTimeFormat('en', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        })

        return dateTimeFormat.format(new Date(val))
    }

    typeToEmoji (type: number): string {
        switch (type) {
            case 0: return '👍'
            case 1: return '👎'
            case 2: return '🔥'
            case 3: return '🤬'
        }

        return '😨'
    }

    countOf (type: number): number {
        return this.feedbacks.filter((it) => { return it.type === type }).length
    }

    created () {

        this.$vs.setColor('background', '#212529')

        const srv = new FeedbackService()
        srv.loadServices().then((res) => {
            this.isLoading = false
            this.feedbacks = res
            console.log(this.feedbacks)
            const dateTimeFormat = new Intl.DateTimeFormat('en', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            })

            console.log()
        }).catch((err) => {
            window.alert(`Error: ${err}`)
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

.date-p {
    font-family: 'Comic Sans MS', 'Comic Sans', cursive;
    font-size: 10pt;
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

</style>
