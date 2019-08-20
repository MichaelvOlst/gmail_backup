<template>
  <v-container>
    <v-form ref="formSettings">
      
      <v-card class="mb-5" v-for="option in settings.storage_options" :key="option.name">
        <v-card-title>
          <span class="mr-2">{{ option.name }}</span>
          <v-divider></v-divider>
           <v-btn class="mx-2 mt-1" @click="save()" fab dark small color="primary">
            <v-icon>save</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text>     

          <v-switch
            v-model="option.active"
            label="Active"
            color="primary"
          ></v-switch>   

          <v-text-field
            v-model="option.path"
            label="Path"
            required
          ></v-text-field>

          <v-text-field v-for="(value, key) in Object.keys(option.config)" :key="key"
            v-model="option.config[value]"
            :label="value"
            required
          ></v-text-field>   

        </v-card-text>
      </v-card>
    </v-form>
  </v-container>
</template>

<script>
  import { mapState } from 'vuex'
  import { GET_SETTINGS, SAVE_SETTINGS } from './../store/modules/types'  

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
        } catch(e) {
          console.log(e);
        }
      }
    }

  }
</script>