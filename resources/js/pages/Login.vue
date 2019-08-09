<template>
  <div>
    <v-card class="elevation-12">
      <v-card-text>
        <v-form @keyup.enter.prevent="login()">
          <v-text-field
            label="Email"
            name="email"
            prepend-icon="person"
            type="text"
            v-model="form.email"
          ></v-text-field>

          <v-text-field
            id="password"
            label="Password"
            name="password"
            prepend-icon="lock"
            type="password"
            v-model="form.password"
          ></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-btn class="ml-2" color="primary" @click.prevent="login()">Login</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { LOGIN, CHECK_AUTH } from './../store/modules/types';

  export default {
    data () {
      return {
        form: {
          email: "",
          password: ""
        }
      }
    },

    created () {
      this.$store.dispatch(CHECK_AUTH)
      .then( () => {
        this.$router.push({ name: "dashboard"})
      })
      .catch(() => { })      
    },

    methods: {
      login() {
        this.$store.dispatch(LOGIN, this.form)
          .then(()=>{
            let redirect = this.$route.query.redirect || "dashboard"
            this.$router.push({ path: redirect})
          })
      }
    }
  }
</script>