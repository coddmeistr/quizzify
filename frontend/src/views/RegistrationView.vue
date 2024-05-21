<template>
  <div class="login-view">
    <v-container fluid fill-height>
      <v-layout class="flex" align-center justify-center>
        <div class="flex">
          <v-card class="elevation-12 body">
            <v-toolbar dark color="primary">
              <v-toolbar-title>Register</v-toolbar-title>
            </v-toolbar>
            <v-card-text>
              <v-form>
                <v-text-field
                    v-model.trim="username"
                    prepend-icon="mdi-account"
                    name="login"
                    label="Login"
                    type="text"
                    :rules="loginFormRules"
                    :error="isLoginFailed"
                ></v-text-field>
                <v-text-field
                    v-model="password"
                    uuid="password"
                    prepend-icon="mdi-lock"
                    name="password"
                    label="Password"
                    type="password"
                    :rules="loginFormRules"
                    :error="isLoginFailed"
                    :error-messages="getErrorMessages"
                    @keypress.enter="tryRegister"
                ></v-text-field>
                <v-text-field
                    v-model="email"
                    uuid="email"
                    prepend-icon="mdi-email"
                    name="email"
                    label="Email"
                    type="email"
                    :rules="loginFormRules"
                    :error="isLoginFailed"
                    :error-messages="getErrorMessages"
                    @keypress.enter="tryRegister"
                ></v-text-field>
              </v-form>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="primary" @click="tryRegister" :loading="loginLoading">
                Login
              </v-btn>
            </v-card-actions>
          </v-card>
        </div>
      </v-layout>
    </v-container>
  </div>
</template>

<script>
import snackbar from "@/mixins/snackbar";

export default {
  name: "register-view",
  mixins: [snackbar],
  data() {
    return {
      loginFormRules: [(v) => !!v || "Required"],
      loginLoading: false,
      isLoginFailed: false,
      username: "",
      password: "",
      email: "",
      type: "Standard",
    };
  },
  methods: {
    tryRegister() {
      this.loginLoading = true;
      (this.isLoginFailed = false),
          this.$store
              .dispatch("auth/register", {
                login: this.username,
                password: this.password,
                email: this.email,
              })
              .then(() => {
              })
              .catch((error) => {
                console.log(error);
                this.showSnackbarError({
                  message: error.response.data.message || "Error during registration",
                });
                if (error.response && error.response.status != 200) {
                  this.isLoginFailed = true;
                }
              })
              .finally(() => {
                this.loginLoading = false;
              });
    },
  },
  computed: {
    getErrorMessages() {
      return this.isLoginFailed ? ["username or password is not correct"] : [];
    },
  },
};
</script>

<style>
.login-view {
  height: 100%;
}

.type-select {
  margin-left: 30px;
}

.flex{
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 50px;
}

.body{
  min-width: 500px;
}
</style>