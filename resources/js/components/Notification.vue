<template>
    <v-snackbar v-model="show" :timeout=2000>
        {{message}}
        <v-btn color="red" @click.native="show = false">
            <v-icon>close</v-icon>
        </v-btn>
    </v-snackbar>
</template>

<script>
    import { NOTIFY } from './../store/modules/types' 

    export default {
        data () {
            return {
                show: false,
                message: ""
            }
        },

        created () {
            this.$store.watch(state => state.notification.message, () => {
            const msg = this.$store.state.notification.message
                if (msg !== '') {
                    this.show = true
                    this.message = msg
                    this.$store.commit(NOTIFY, '')
                }
            })
        }
    }
</script>