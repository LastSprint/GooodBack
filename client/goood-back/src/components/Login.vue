<template>
    <div class="pageCenter">
        <div class="center grid">
          <vs-row>
            <vs-col vs-type="flex" vs-justify="center" vs-align="center" w="12">
                <div id="container">
                        <div class="flip-card" @click="flip">
                          <div class="flip-card-inner" id="flip-it-elem">
                            <div class="flip-card-front">
                              <img src="../assets/login_logo.png" width="118px" height="120px" class="grow crossRotate" id="gopher-logo">
                            </div>
                            <div class="flip-card-back" id="container">
                              <img src="../assets/google_logo.png" width="90px" height="90px">
                            </div>
                          </div>
                        </div>
                  <p>GoodBack</p>
                </div>
            </vs-col>
           </vs-row>
          <vs-row>
            <vs-col vs-type="flex" vs-justify="center" vs-align="center" w="12">
                <div id="container">
                  <vs-button
                    size="xl"
                    circle
                    color="warn"
                    gradient    
                    to="/feedback-form"
                  >
                    📬 Send Feedback
                  </vs-button>
                </div>
            </vs-col>
          </vs-row>
        </div>
    </div>
</template>

<script lang="ts">

import { Component, Prop, Vue } from 'vue-property-decorator'
import AuthService from '../services/AuthService'
import * as global from '../common/GlobalConstants'
import 'vuesax/dist/vuesax.css'

@Component
  export default class Login extends Vue {
    clientId!: string;
    _isFlipped = false;
    // service: AuthService;

    handleSignIn (): void {
      const clientId = global.GOOGLE_CLIENT_ID
      const responseType = 'code'
      const scope = 'https://www.googleapis.com/auth/userinfo.email'
      const redirectUri = global.REDIRECT_URL
      const url = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}&response_type=${responseType}&scope=${scope}&redirect_uri=${redirectUri}`
      console.log(url)
      window.location.href = url
    }

    created () {
      const srv = new AuthService()
      if (srv.isAuthorizaed()) {
        // this.$router.replace({ name: 'Feedback' })
      }
    }

    flip () {

      if (this._isFlipped) {
        this.handleSignIn()
        return
      }

      this._isFlipped = true

      var elem = document.getElementById('flip-it-elem')
      elem!.style.transform = 'rotateY(180deg)'

      var gopher = document.getElementById('gopher-logo')
      gopher!.style.transform = 'scale(0.9)'

      setInterval(() => { 
        elem!.style.transform = 'rotateY(0deg)'
        this._isFlipped = false 
        gopher!.style.transform = ''
        }, 2500)
    }
}
</script>

<style scoped>

.pageCenter {
    position: fixed;
    top: 50%;
    left: 50%;
    -webkit-transform: translate(-50%, -50%);
    transform: translate(-50%, -50%);
}

.vcontainer {
    min-height: 10em;
    display: table-cell;
    vertical-align: middle;
}

#container {
  display: flex;
  justify-content: center;
  align-items: center;
}

p{
 font-family: 'Comic Sans MS', 'Comic Sans', cursive;
 font-size: 65pt;
 margin-left: 16pt;
 color: #e9ecef;
}

.grow { transition: all .2s ease-in-out; }
.grow:hover { transform: scale(1.1); }

/* The flip card container - set the width and height to whatever you want. We have added the border property to demonstrate that the flip itself goes out of the box on hover (remove perspective if you don't want the 3D effect */
.flip-card {
  background-color: transparent;
  width: 120px;
  height: 120px;
  perspective: 1000px; /* Remove this if you don't want the 3D effect */
}

/* This container is needed to position the front and back side */
.flip-card-inner {
  position: relative;
  width: 100%;
  height: 100%;
  text-align: center;
  transition: transform 0.8s;
  transform-style: preserve-3d;
}

/* Do an horizontal flip when you move the mouse over the flip box container
.flip-card:hover .flip-card-inner {
  transform: rotateY(180deg);
} */

/* Position the front and back side */
.flip-card-front, .flip-card-back {
  position: absolute;
  width: 100%;
  height: 100%;
  -webkit-backface-visibility: hidden; /* Safari */
  backface-visibility: hidden;
}

/* Style the front side (fallback if image is missing) */
.flip-card-front {
  background-color: transparent;
}

/* Style the back side */
.flip-card-back {
  background-color: whitesmoke;
  transform: rotateY(180deg);
  border-radius: 100%;
}

</style>
