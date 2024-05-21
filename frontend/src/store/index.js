
import Vuex from "vuex";

// modules
import auth from "./auth";
import tests from "./tests";


const store = new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    auth,
    tests,
  },
});

export default store;

export const useStore = () => store;
