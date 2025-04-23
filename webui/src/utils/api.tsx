import axios from "axios";

axios.defaults.baseURL = "/api/v1/";

export default class API {
    static handler = function () {
        return axios.create({
            baseURL: "/api/v1/",
        })
    };

    static config = {
        get: function () {
            return API.handler().get('/config');
        },
        get_version: function () {
            return API.handler().get('/config/version');
        },
    };

    static jobs = {
        list: async function (p = {}) {
            return API.handler().get('/jobs', { params: p });
        },
        get: function (uuid: string) {
            return API.handler().get('/jobs/' + uuid);
        },
        getLogs: function (uuid: string, backup_path: string) {
            let params = { "bp": backup_path };
            let sp = new URLSearchParams(params)
            return API.handler().get('/jobs/' + uuid + '/logs?' + sp.toString());
        },
    };

    static clients = {
        list: async function (p = {}) {
            return API.handler().get('/clients', p);
        },
        get: function (name: string) {
            return API.handler().get('/clients/' + name);
        },
        ping: function (name: string) {
            return API.handler().get('/clients/' + name + "/ping");
        },
    };

    static modules = {
        list: async function (p = {}) {
            return API.handler().get('/modules', p);
        },
        get: function (name: string) {
            return API.handler().get('/modules/' + name);
        },
    };

    static images = {
        list: async function (p = {}) {
            return API.handler().get('/images', { params: p });
        },
        get: function (uuid: string) {
            return API.handler().get('/images/' + uuid);
        },
        stats: function () {
            return API.handler().get('/images/stats');
        },
    };

    static repos = {
        list: async function (p = {}) {
            return API.handler().get('/repositories', { params: p });
        },
        get: function (uuid: string) {
            return API.handler().get('/repositories/' + uuid);
        },
    };

}