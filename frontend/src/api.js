import {BackendHost} from "@/consts";
import axios from "axios";
import {useStore} from "@/store";

// Add a response interceptor
axios.interceptors.response.use(function (response) {
    // Any status code that lie within the range of 2xx cause this function to trigger
    // Do something with response data
    return response;
}, function (error) {
    // Any status codes that falls outside the range of 2xx cause this function to trigger
    // Do something with response error
    return Promise.reject(error);
});

export function getAxios(){
    return axios.create({
        baseURL: "http://"+BackendHost
    });
}

export function getAuthConfig(){
    let store = useStore();
    let userdata = store.getters["auth/userdata"];
    if (!userdata || !userdata.userId) return {};
    return {
        headers: {
            'Auth-User-Info': JSON.stringify({id: parseInt(userdata.userId), permissions: userdata.permissions})
        }
    }
}