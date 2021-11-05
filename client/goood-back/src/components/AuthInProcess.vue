<template>
    <div>
        I'M IN PROCESS. FUCK OFF
    </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import AuthService from '../services/AuthService'
import VueRouter from 'vue-router'

@Component
export default class AuthInProcess extends Vue {
    service!: AuthService

    created () {
        this.service = new AuthService()
        console.log(this.$router.currentRoute.query)
        this.service.authorize(this.$router.currentRoute.query.code as string)
            .then((res) => {
                this.$router.replace({ name: 'Feedback' })
            }).catch((err) => {
                window.alert(`Error ${err.response.data}`)

                this.$router.replace({ name: 'Login' })
            })
    }
}
</script>
