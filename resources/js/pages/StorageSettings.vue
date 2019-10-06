<template>
  <v-container>

    <!-- {{ settings }} -->

    <v-form ref="formSettings">
      <div  v-if="settings.storage_options">
        <v-card class="mb-5" v-for="(provider) in Object.keys(settings.storage_options)" :key="provider">

          <!-- {{ settings.storage_options[provider].Config }} -->
          
          <v-card-title>
            <span class="mr-2">{{ settings.storage_options[provider].StorageOption.name }}</span>
            <v-divider></v-divider>
            <v-btn class="mx-2 mt-1" @click="save()" fab dark small color="primary">
              <v-icon>save</v-icon>
            </v-btn>
          </v-card-title>
          <v-card-text>     

            <v-switch
              v-model="settings.storage_options[provider].StorageOption.active"
              label="Active"
              color="primary"
            ></v-switch>   

            <v-text-field v-for="(value, key) in Object.keys(settings.storage_options[provider].Config)" :key="key"
              v-model="settings.storage_options[provider].Config[value]"
              :label="value"
              required
            ></v-text-field>   

          </v-card-text>
        </v-card>
      </div>
    </v-form>
  </v-container>
</template>

<script>
  import { mapState } from 'vuex'
  import { GET_SETTINGS, SAVE_SETTINGS, NOTIFY } from './../store/modules/types' 

  export default {

    async created () {
      try {
        let response = this.$store.dispatch(GET_SETTINGS)  
      } catch (e) {
        console.log(e);
      }
    },

    computed: {
      ...mapState({
        settings: state => state.settings.settings,
      }),
    },

    methods: {
      async save() {
        try {
          let response = this.$store.dispatch(SAVE_SETTINGS, this.settings)
          this.$store.commit(NOTIFY, "Settings saved")
        } catch(e) {
          console.log(e);
        }
      }
    }

  }
</script>