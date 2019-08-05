<template>
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
                    <v-text-field v-model="form.email" label="Emailaddress"></v-text-field>
                  </v-flex>
                  <v-flex xs12 sm12 md12>
                    <!-- 4/nQH0XkilbxhMtkm6z2fqi7lyB2I5zd7Cxq4U4kwhD5IWF8uNiPDrH-E -->
                    <v-text-field v-model="form.encryption_key" label="Encryption key"></v-text-field>
                      <a href="#" @click.prevent="generateKey()" target="blank">Generate key</a>

                  </v-flex>
                  <v-flex xs12 sm12 md12>
                    <v-text-field v-model="form.accesstoken" label="Google accesstoken"></v-text-field>
                    <a v-if="getAccessTokenURL" :href="getAccessTokenURL" target="blank">Get accesstoken</a>
                  </v-flex>
                  <v-flex xs12 sm12 md12>
                    <v-switch v-model="form.attachments" label="Attachments"></v-switch>
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
    <template v-slot:item.action="{ item }">
      <v-icon small class="mr-2" @click="editItem(item)">edit</v-icon>
      <v-icon small @click="deleteItem(item)">delete</v-icon>
    </template>
  </v-data-table>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { GOOGLE_URL } from './../store/modules/types'  

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
        { text: 'Last date', value: 'backup_date' },
        { text: 'percentage', value: 'percentage' },
        { text: 'Actions', value: 'action', sortable: false },
      ],
      accounts: [],
      editedIndex: -1,
      form: {
        email: '',
        encryption_key: '',
        attachments: true,
        accesstoken: '',
      }
    }),

    computed: {
      formTitle () {
        return this.editedIndex === -1 ? 'New Account' : 'Edit Account'
      },
      ...mapGetters([
        'getAccessTokenURL',
      ])
    },

    watch: {
      dialog (val) {
        val || this.close()
      },
    },

    created () {
      this.$store.dispatch(GOOGLE_URL).then((data) => {
        
      })
      .catch(()=> {
        console.error("could not retrieve URL for google auth")
      })

    },

    methods: {

      generateKey() {
        this.form.encryption_key =  Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
      },

      editItem (item) {
        this.editedIndex = this.desserts.indexOf(item)
        this.editedItem = Object.assign({}, item)
        this.dialog = true
      },

      deleteItem (item) {
        const index = this.desserts.indexOf(item)
        confirm('Are you sure you want to delete this item?') && this.desserts.splice(index, 1)
      },

      close () {
        this.dialog = false
        setTimeout(() => {
          this.editedItem = Object.assign({}, this.defaultItem)
          this.editedIndex = -1
        }, 300)
      },

      save () {
        if (this.editedIndex > -1) {
          Object.assign(this.desserts[this.editedIndex], this.editedItem)
        } else {
          this.desserts.push(this.editedItem)
        }
        this.close()
      },
    },
  }
</script>