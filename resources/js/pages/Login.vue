<template>
  <div>
    <v-card class="elevation-12">
      <v-card-text>
        <v-form>
          <v-text-field
            label="Login"
            name="login"
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
        <v-spacer></v-spacer>
        <v-btn color="primary" @click.prevent="login()">Login</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
  import { LOGIN } from './../store/modules/types';

  export default {
    data () {
      return {
        form: {
          email: "",
          password: ""
        }
      }
    },

    methods: {
      login() {
        this.$store
        .dispatch(LOGIN, this.form)
        .then(() => this.$router.push({ name: "dashboard" }));
      }
    }
  }
</script>