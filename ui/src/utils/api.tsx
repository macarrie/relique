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
    };

    static clients = {
        list: function (p = {}) {
            return API.handler().get('/clients', p);
        },
        ssh_ping: function (id :number) {
            return API.handler().post("/clients/" + id + "/ssh_ping");
        },
    };
}

