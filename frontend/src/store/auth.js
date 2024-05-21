import {getAxios} from "@/api";
import Cookies from "js-cookie";
import router from "@/router";

const COOKIES_NAME = "quizzify-token";

export default {
  namespaced: true,
  state: {
    token: "",
    userdata: {},
    appURL: "",
    accounts: []
  },
  mutations: {
    setToken(state, token) {
      const expires = new Date(Date.now() + 7776e6);

      state.token = token;
      Cookies.set(COOKIES_NAME, token, { expires });
      // document.cookie = `${COOKIES_NAME}=${token}; path=/; expires=${expires}`;
    },
    setUserdata(state, data) {
      state.userdata = data;
    },
    setAppURL(state, data) {
      state.appURL = data;
    },
    setAccounts(state, data) {
      state.accounts = data;
    }
  },
  actions: {
    deleteAccount(_, id) {
      return new Promise((resolve, reject) => {
        getAxios().delete("/sso/account?id="+id)
            .then((response) => {
              resolve(response);
            })
            .catch((error) => {
              reject(error);
            });
      });
    },
    accountsList({ commit }) {
      return new Promise((resolve, reject) => {
        getAxios().get("/sso/accounts")
            .then((response) => {
              commit("setAccounts", response.data.accounts);
              resolve(response);
            })
            .catch((error) => {
              reject(error);
            });
      });
    },
    login({ commit }, { login, password }) {
      return new Promise((resolve, reject) => {
        getAxios().post("/sso/login", {
            login: login,
            password: password,
            app_id: 1,
        })
            .then((response) => {
              commit("setToken", response.data.token);
              resolve(response);
            })
            .catch((error) => {
              reject(error);
            });
      });
    },
    register(_, { login, password, email }) {
      return new Promise((resolve, reject) => {
        getAxios().post("/sso/register", {
          login: login,
          password: password,
          email: email
        })
            .then((response) => {
              router.push({ name: "Login" });
              resolve(response);
            })
            .catch((error) => {
              reject(error);
            });
      });
    },
    logout({ commit }) {
      commit("setToken", "");
      commit("setUserdata", {});
      Cookies.remove(COOKIES_NAME);
      router.push({ name: "Login" });
    },

    fetchUserData({commit}, {token}) {
      return new Promise((resolve, reject) => {
          if (token === undefined || token === "") {
              token = Cookies.get(COOKIES_NAME);
              if (token !== undefined && token !== "") {
                  commit("setToken", token)
              } else {
                  reject("no token")
                  return
              }
          }
        let config = {
          headers: {
            Authorization: `Bearer `+token,
          }
        }
        getAxios().get(`/sso/account?token=${token}`, config)
            .then((response) => {
              commit("setUserdata", response.data);
              resolve(response);
            })
            .catch((error) => {
              reject(error);
            });
      });
    },
  },
  getters: {
    isLoggedIn(state) {
      return state.token.length > 0;
    },
    userdata(state) {
      return state.userdata;
    },
    token(state) {
      return state.token;
    },
    accounts(state) {
      return state.accounts;
    }
  },
};
