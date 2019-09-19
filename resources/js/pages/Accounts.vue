<template>


  <div>
     <!-- {{ accounts }} -->

     <!-- {{ settings.active_storage_options }} -->
  <v-data-table
    :headers="headers"
    :items="accounts"
    no-data-text="No accounts have been created yet"
    class="elevation-1"
  >

    <template v-slot:top>
      <v-toolbar flat color="white">
        <v-dialog v-model="dialog" max-width="1000px">
          <template v-slot:activator="{ on }">
            <v-btn color="primary" dark class="mb-2" v-on="on">New account</v-btn>
          </template>
          <v-card>
            <v-card-title>
              <span class="headline">{{ formTitle }}</span>
            </v-card-title>

            <v-card-text>
              <v-container grid-list-md>
                <v-layout wrap>
                  <v-flex xs12 sm12 md12>
                    <v-text-field v-model="form.email" :error="errors.email !== null" :error-messages="errors.email" label="Emailaddress" required></v-text-field>
                  </v-flex>
                  <v-flex xs12 sm12 md12>
                    <v-text-field v-model="form.encryption_key" :error="errors.encryption_key !== null" :error-messages="errors.encryption_key" label="Encryption key" required></v-text-field>
                    <a href="#" @click.prevent="generateKey()" target="blank">Generate key</a>
                  </v-flex>
                  <v-flex xs12 sm12 md12>
                    <v-text-field v-model="form.google_token" :error="errors.google_token !== null" :error-messages="errors.google_token" label="Google token" required></v-text-field>
                    <a v-if="getTokenURL" :href="getTokenURL" target="blank">Get token</a>
                  </v-flex>

                  <v-flex xs12 sm12 md12>
                    <v-text-field v-model="form.upload_path" :error="errors.upload_path !== null" :error-messages="errors.upload_path" label="Upload path" required></v-text-field>
                  </v-flex>

                  <v-flex xs12 sm12 md12>
                    <v-select
                      :items="settings.active_storage_options"
                      item-text="name"
                      item-value="option"
                      label="Storage option"
                      v-model="form.storage_provider"
                      :error="errors.upload_path !== null"
                      :error-messages="errors.upload_path" 
                    ></v-select>
                  </v-flex>                  
                </v-layout>
              </v-container>  
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" text @click="close">Cancel</v-btn>
              <v-btn color="blue darken-1" text @click="save">Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-toolbar>
    </template>

    <template v-slot:item.storage_provider="{ item }">
      <div>{{ getActiveStorage(item) }} </div>
    </template>

    <template v-slot:item.backup="{ item }">
      <v-layout justify-space-between>
        <v-icon small @click="backup(item)">backup</v-icon>{{ item.backup_progress_message }}
      </v-layout>
    </template>

    <template v-slot:item.action="{ item }">
      <v-icon small class="mr-2" @click="editItem(item)">edit</v-icon>
      <v-icon small @click="deleteItem(item)">delete</v-icon>
    </template>
  </v-data-table>
  </div>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import { GOOGLE_URL, SAVE_ACCOUNT, ALL_ACCOUNTS, GET_ACCOUNT, DELETE_ACCOUNT, BACKUP_ACCOUNT, GET_SETTINGS } from './../store/modules/types'  
  import { log } from 'util';

  export default {
    data: () => ({
      dialog: false,
      headers: [
        {
          text: 'Email',
          align: 'left',
          sortable: true,
          value: 'email',
        },
        { text: 'Storage', value: 'storage_provider' },
        { text: 'Last date', value: 'backup_date' },
        { text: 'backup', value: 'backup' },
        { text: 'Actions', value: 'action', sortable: false },
      ],
      form: {
        email: '',
        encryption_key: '',
        attachments: true,
        google_token: '',
        storage_provider: null,
      },
      errors: {
        email: null,
        encryption_key: null,
        google_token: null,
        upload_path: null,
        storage_provider: null,
      },
    }),

    computed: {
      formTitle () {
        return this.form.id ? 'Edit Account' : 'New Account'
      },
      ...mapGetters([
        'getTokenURL',
      ]),
      ...mapState({
        accounts: state => state.accounts.accounts,
        settings: state => state.settings,
      }),
    },


    watch: {
      dialog (val) {
        val || this.close()
      },

      form: {
        handler (val) {
          this.errors.email = null
          this.errors.encryption_key = null
          this.errors.google_token = null
        },
        deep: true
      }
    },

    async created () {

      try {
        let response = this.$store.dispatch(GOOGLE_URL)  
      } catch (e) {
        console.log(e);
      }

      try {
        let response = this.$store.dispatch(GET_SETTINGS)  
      } catch (e) {
        console.log(e);
      }

      this.getAllAccounts()
    },

    methods: {

      getActiveStorage (item) {       

        const options = this.settings.active_storage_options
        let providerName = null
        for (let index = 0; index < options.length; index++) {
          const provider = options[index];
          if (provider.option == item.storage_provider) {
            providerName = provider.name
          }
        }
        return providerName || "Unknown"
      },

      generateKey() {
        this.form.encryption_key = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
      },

       async backup (item) {
        if(!confirm("Are you sure?")) {
          return
        }

        let websocket = new WebSocket(`ws:/${window.location.host}/api/backup/${item.id}`);
        console.log("Attempting Connection...");

        websocket.onopen = function(event) {
          console.log("Successfully connected to websocket server");
        };

        websocket.onerror = function(error) {
          console.log("Error connecting to websocket server");
          console.log(error);
          websocket.close();
        };

        websocket.onmessage = function(event) {
          // parse the event data sent from our websocket server
          let data = JSON.parse(event.data);

          if (data.error) {
            alert(`Error occured: ${data.error}`)
            item.backup_progress_message = data.error
            websocket.close()
            return
          }

          if(data.backup_progress_message == "done") {
            alert("Done")
            item.backup_progress_message = data.backup_progress_message
            websocket.close()
            return
          }

          item.backup_progress_message = data.backup_progress_message

          // populate our `sub` element with the total subscriber counter for our
          // channel
          // console.log(data)
        };

        // try {
        //   let response = await this.$store.dispatch(BACKUP_ACCOUNT, item.id)
        //   console.log(response)
        // } catch (e) {
        //   console.log(e)
        // }
      },

      async getAllAccounts() {
        try {
          let response = this.$store.dispatch(ALL_ACCOUNTS)
        } catch (e) {
          console.log(e);
        }
      },

      async editItem (item) {
        try {
          let response = await this.$store.dispatch(GET_ACCOUNT, item.id)
          this.form = response;
          this.dialog = true
        } catch (e) {
          console.log(e.error);          
        }
      },

      async deleteItem (item) {
        if(!confirm("Are you sure?")) {
          return
        }

        try {
          let response = await this.$store.dispatch(DELETE_ACCOUNT, item.id)
          this.getAllAccounts()
        } catch (e) {

        }
      },

      close () {
        this.dialog = false
        this.form = {
          email: '',
          encryption_key: '',
          attachments: true,
          google_token: '',
          storage_provider: null,
        }
      },

      async save () {
        try {
          let response = await this.$store.dispatch(SAVE_ACCOUNT, this.form)
          // this.form = {
          //   email: '',
          //   encryption_key: '',
          //   attachments: true,
          //   google_token: '',
          //   storage_provider: null,
          // }
          this.close();
          this.getAllAccounts()
        } catch(err) {          
          this.errors = err
        }
      },
    },
  }
</script>