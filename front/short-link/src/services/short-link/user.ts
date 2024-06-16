import {request} from "@umijs/max";


export async function login(data) {
    return request("http://127.0.0.1:8088/api/v1/admin/users/login",{
        method: 'POST',
        data
    })
}


export async function register(data) {
    return request("http://127.0.0.1:8088/api/v1/admin/users/register", {
        method: 'POST',
        data,
    });
}

export async function queryCurrentUser(){
    return request("http://127.0.0.1:8088/api/v1/admin/users/current-user", {
        method: 'GET',
    });
}
