import {getAxios, getAuthConfig,getUserData} from "@/api";

export default {
    namespaced: true,
    state: {
        tests: [],
        results: [],
        test: {}
    },
    mutations: {
        setTests(state, data) {
            state.tests = data
        },
        setResults(state, data) {
            state.results = data;
        },
        setTest(state, data) {
            state.test = data;
        },
    },
    actions: {
        getTest({commit}, id) {
            try {
                return new Promise((resolve, reject) => {
                    getAxios()
                        .get("/api/tests/" + id, getAuthConfig())
                        .then((response) => {
                            commit("setTest", response.data.payload)
                            resolve(response)
                        })
                        .catch((error) => {
                            reject(error)
                        })
                });
            } catch (e) {
                console.log(e)
            }
        },
        deleteTest(_, id) {
            return new Promise((resolve, reject) => {
                getAxios()
                    .delete("/api/tests/" + id, getAuthConfig())
                    .then((response) => {
                        resolve(response)
                    })
                    .catch((error) => {
                        reject(error)
                    })
            });
        },
        getTests({commit}, withAnswers) {
            return new Promise((resolve, reject) => {
                getAxios()
                    .get("/api/tests?withAnswers=" + withAnswers)
                    .then((response) => {
                        commit("setTests", response.data.payload)
                        resolve(response)
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        getResults({commit}) {
            return new Promise((resolve, reject) => {
                getAxios()
                    .get("/api/tests/results")
                    .then((response) => {
                        commit("setResults", response.data.payload)
                        resolve(response)
                    })
                    .catch((error) => {
                        reject(error)
                    })
            });
        },
        sendResult(_, {test_id, answers}) {
            return new Promise((resolve, reject) => {
                getAxios()
                    .post(`/api/tests/${test_id}/apply`, {user_answers: answers}, getAuthConfig())
                    .then((response) => {
                        resolve(response)
                        alert("Результат сохранен")
                    })
                    .catch((error) => {
                        alert("Ошибка сохранения результата")
                        reject(error)
                    })
            });
        },
        createTest(_, {test}) {
            return new Promise((resolve, reject) => {
                let userdata = getUserData();
                test.creator_id = parseInt(userdata.userId)
                getAxios()
                    .post(`/api/tests`, {...test}, getAuthConfig())
                    .then((response) => {
                        resolve(response)
                        alert("Тест успешно создан")
                    })
                    .catch((error) => {
                        alert("Ошибка создания теста")
                        reject(error)
                    })
            });
        },
    },
    getters: {
        tests(state) {
            return state.tests;
        },
        results(state) {
            return state.results;
        },
        test(state) {
            return state.test;
        }
    }
}
