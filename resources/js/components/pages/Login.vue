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
          ></v-text-field>

          <v-text-field
            id="password"
            label="Password"
            name="password"
            prepend-icon="lock"
            type="password"
          ></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary">Login</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
  export default {
    
  }
</script>