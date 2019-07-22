<template>
  <v-app id="gmail_backup">
    <v-navigation-drawer v-model="drawer" fixed clipped class="grey lighten-4" app>
      <v-list dense class="grey lighten-4">
        <template >
             <v-list-tile>
                <v-list-tile-action>
                <v-icon>dashboard</v-icon>
                </v-list-tile-action>
                <v-list-tile-title>Dashboard</v-list-tile-title>
            </v-list-tile>
            <v-list-tile>
                <v-list-tile-action>
                <v-icon>people</v-icon>
                </v-list-tile-action>
                <v-list-tile-title>Account</v-list-tile-title>
            </v-list-tile>
             <v-list-tile>
                <v-list-tile-action>
                <v-icon>settings</v-icon>
                </v-list-tile-action>
                <v-list-tile-title>Settings</v-list-tile-title>
            </v-list-tile>
        </template>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar color="primary" app clipped-left>
      <v-toolbar-side-icon @click="drawer = !drawer"></v-toolbar-side-icon>
      <span class="title ml-3 mr-5">Gmail Backup</span>
    </v-toolbar>
    <v-content>
      <v-container fluid fill-height class="grey lighten-4">
        <v-layout justify-center align-center>
         
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
  export default {
    
  }
</script>

<style>
  #keep .v-navigation-drawer__border{ display: none;}
</style>