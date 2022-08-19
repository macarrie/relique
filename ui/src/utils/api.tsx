import axios from "axios";

axios.defaults.baseURL = "/api/v1/";

export default class API {
    static handler = function() {
        return axios.create({
            baseURL: "/api/v1/",
            //headers: {'Authorization': 'Bearer ' + Auth.getToken()}
        })
    };

    static jobs = {
        list: function (p = {}) {
            return API.handler().post('/jobs', p);
        },
        get: function (uuid :string) {
            return API.handler().get('/jobs/' +uuid);
        },
    };

    static clients = {
        list: function (p = {}) {
            return API.handler().get('/clients', p);
        },
        get: function (id :number) {
            return API.handler().get('/clients/' +id);
        },
    };

    static modules = {
        list: function (p = {}) {
            return API.handler().get('/modules', p);
        },
        get: function (name :string) {
            return API.handler().get('/modules/' +name);
        },
    };
}

