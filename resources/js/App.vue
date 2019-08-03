<template>
  <v-app id="gmail_backup">
    <v-navigation-drawer app clipped>
       <v-list dense class="grey lighten-4">        
          <v-list-item to="dashboard">
            <v-list-item-action>
            <v-icon>dashboard</v-icon>
            </v-list-item-action>
            <v-list-item-title>
              Dashboard
            </v-list-item-title>
        </v-list-item>
        <v-list-item to="accounts">
            <v-list-item-action>
            <v-icon>people</v-icon>
            </v-list-item-action>
            <v-list-item-title>Account</v-list-item-title>
        </v-list-item>
        <v-list-item to="settings">
            <v-list-item-action>
            <v-icon>settings</v-icon>
            </v-list-item-action>
            <v-list-item-title>Settings</v-list-item-title>
        </v-list-item>
        <v-list-item v-if="isAuthenticated" @click.prevent="logout()">
            <v-list-item-action>
            <v-icon>settings</v-icon>
            </v-list-item-action>
            <v-list-item-title>Logout</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left dense color="primary">
      <v-toolbar-title class="mr-12 align-center">
      <span class="title ml-3 mr-5">Gmail Backup</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-layout align-center style="max-width: 650px"></v-layout>
    </v-app-bar>

    <v-content>
      <v-container fill-height>
        <v-layout justify-center align-center>
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
