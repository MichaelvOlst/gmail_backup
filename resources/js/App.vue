<template>
  <v-app id="gmail_backup">
    <v-navigation-drawer fixed clipped class="grey lighten-4" app>
      <v-list dense class="grey lighten-4">        
          <v-list-tile to="dashboard">
            <v-list-tile-action>
            <v-icon>dashboard</v-icon>
            </v-list-tile-action>
            <v-list-tile-title>
              Dashboard
            </v-list-tile-title>
        </v-list-tile>
        <v-list-tile to="accounts">
            <v-list-tile-action>
            <v-icon>people</v-icon>
            </v-list-tile-action>
            <v-list-tile-title>Account</v-list-tile-title>
        </v-list-tile>
        <v-list-tile to="settings">
            <v-list-tile-action>
            <v-icon>settings</v-icon>
            </v-list-tile-action>
            <v-list-tile-title>Settings</v-list-tile-title>
        </v-list-tile>
        <v-list-tile v-if="isAuthenticated" @click.prevent="logout()">
            <v-list-tile-action>
            <v-icon>settings</v-icon>
            </v-list-tile-action>
            <v-list-tile-title>Logout</v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar color="primary" app clipped-left>
      <span class="title ml-3 mr-5">Gmail Backup</span>
    </v-toolbar>
    <v-content>
      <v-container fluid fill-height>  
        <v-layout justify-center align-center >
          <v-flex md12 lg12 class="pa-0 ma-0 ">
            <router-view></router-view>
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'
  import { LOGOUT } from './store/modules/types'  

  export default {
    name: "App",

    computed: {
      // mix the getters into computed with object spread operator
      ...mapGetters([
        'isAuthenticated',
      ])
    },

    methods: {
     logout() {
        this.$store.dispatch(LOGOUT).then(() => {
          this.$router.push({ name: "login" });
        });
      }

    }
  }
</script>
